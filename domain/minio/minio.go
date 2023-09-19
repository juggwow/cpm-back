package minio

import (
	"context"
	"cpm-rad-backend/domain/auth"
	"cpm-rad-backend/domain/config"
	"cpm-rad-backend/domain/connection"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/utils"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/inhies/go-bytesize"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

var (
	minioClient = &minio.Client{}
	// ctx         = context.Background()
)

type Configuration struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
	BucketName string
}

type client struct {
	client *minio.Client
}

type DbRadPdfSign struct {
	ID         uint            `gorm:"column:ID"`
	ReportID   uint            `gorm:"column:RAD_ID"`
	Name       string          `gorm:"column:FILE_NAME"`
	Size       decimal.Decimal `gorm:"column:FILE_SIZE"`
	Unit       string          `gorm:"column:FILE_UNIT"`
	Path       string          `gorm:"column:FILE_PATH"`
	CreateBy   string          `gorm:"column:CREATED_BY"`
	UpdateBy   string          `gorm:"column:UPDATED_BY"`
	UpdateDate *time.Time      `gorm:"column:UPDATED_DATE"`
	DelFlag    string          `gorm:"column:DEL_FLAG"`
}

func (DbRadPdfSign) TableName() string {
	return "CPM.RAD_FILE_SIGNATURE"
}

func NewConnection(config Configuration) (err error) {
	minioClient, err = minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return err
	}
	// log.Printf("%#v\n", minioClient)
	return nil
}
func GetClient() Client {
	return &client{
		client: minioClient,
	}
}

type Client interface {
	Upload(ctx context.Context, file *multipart.FileHeader, floder string) (*minio.UploadInfo, string, error)
	Delete(ctx context.Context, filename string, itemID uint) error
	Download(ctx context.Context, filename string) (*minio.Object, string, error)
}

func UploadHandler(db *connection.DBConnection) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)
		claims, _ := auth.GetAuthorizedClaims(c)
		log.Info(strings.Join([]string{claims.EmployeeID, claims.FirstName, claims.LastName}, " "))
		reportid := c.Param("reportid")

		floder := fmt.Sprintf("%s/%s", c.Param("itemid"), reportid)

		form, _ := c.MultipartForm()
		files := form.File["upload"]

		if err := utils.IsValidPdf(files); err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, utils.ReaponseError{Error: err.Error()})
		}

		for _, file := range files {
			b := bytesize.New(float64(file.Size))
			displaySize := b.Format("%.2f ", "", false)
			words := strings.Fields(displaySize)
			size, _ := decimal.NewFromString(words[0])

			fileName := xid.New().String() + "_" + filepath.Clean(file.Filename)
			objectName := fmt.Sprintf("%s/%s", floder, fileName)

			data := DbRadPdfSign{
				ReportID: utils.StringToUint(reportid),
				Name:     file.Filename,
				Size:     size,
				Unit:     words[1],
				Path:     objectName,
				CreateBy: claims.EmployeeID,
			}

			if err := db.CPM.Omit("UpdateBy", "UpdateDate", "DelFlag").Create(&data).Error; err != nil {
				return err
			}

			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			info, err := minioClient.PutObject(c.Request().Context(),
				config.StorageBucketName,
				objectName,
				src,
				-1,
				minio.PutObjectOptions{})

			if err != nil {
				return c.JSON(http.StatusInternalServerError, utils.ReaponseError{Error: err.Error()})
			}

			log.Info(fmt.Sprintf("%#v\n", info))
		}

		return c.JSON(http.StatusCreated, utils.Reaponse{Msg: "Upload File Success"})
	}
}

func (m *client) Upload(ctx context.Context, file *multipart.FileHeader, floder string) (*minio.UploadInfo, string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, "", err
	}
	defer src.Close()

	objectName := xid.New().String() + "_" + filepath.Clean(file.Filename)

	info, err := minioClient.PutObject(ctx,
		config.StorageBucketName,
		fmt.Sprintf("%s/%s", floder, objectName),
		src,
		-1,
		minio.PutObjectOptions{})
	if err != nil {
		return &info, objectName, err
	}

	log.Printf("%#v\n", info)

	return &info, objectName, err
}

func DownloadHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Unwrap(c)

		objectName := c.Param("itemid") + "/" + c.Param("reportid") + "/" + c.Param("filename")
		ext := filepath.Ext(objectName)

		obj, err := minioClient.GetObject(
			c.Request().Context(),
			config.StorageBucketName,
			objectName,
			minio.GetObjectOptions{},
		)
		if err != nil {
			fmt.Println(err)
			return c.NoContent(500)
		}
		defer obj.Close()
		return c.Stream(http.StatusOK, getContentType(ext), obj)
	}
}

func (m *client) Download(ctx context.Context, filename string) (*minio.Object, string, error) {

	ext := filepath.Ext(filename)

	obj, err := minioClient.GetObject(
		ctx,
		config.StorageBucketName,
		filename,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, "", err
	}

	// defer obj.Close()
	// return c.Stream(http.StatusOK, getContentType(ext), obj)
	return obj, getContentType(ext), err
}

func (m *client) Delete(ctx context.Context, objectName string, itemID uint) error {
	return minioClient.RemoveObject(ctx, config.StorageBucketName, objectName, minio.RemoveObjectOptions{})
	// cfags2e44nsipt7ajr40_322363807.jpg
	// cfagql644nsipt7ajr3g_322363807.jpg
	// cfah4cu44nsi942tvdrg_322363807.jpg
}

func getContentType(ext string) string {
	lowerExt := strings.ToLower(ext)
	fmt.Println("lowerExt" + lowerExt)
	if lowerExt == ".png" {
		return "image/png"
	}

	if lowerExt == ".jpg" || lowerExt == ".jpeg" {
		return "image/jpeg"
	}

	if lowerExt == ".pdf" {
		return "application/pdf"
	}

	return "application/octet-stream"
}
