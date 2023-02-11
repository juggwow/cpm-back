package minio

import (
	"context"
	"cpm-rad-backend/domain/config"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/xid"
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
	Upload(ctx context.Context, file *multipart.FileHeader, itemID uint) (*minio.UploadInfo, string, error)
	Delete(ctx context.Context, filename string, itemID uint) error
	Download(ctx context.Context, filename string) (*minio.Object, string, error)
}

func (m *client) Upload(ctx context.Context, file *multipart.FileHeader, itemID uint) (*minio.UploadInfo, string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, "", err
	}
	defer src.Close()

	objectName := xid.New().String() + "_" + filepath.Clean(file.Filename)

	info, err := minioClient.PutObject(ctx,
		config.StorageBucketName,
		fmt.Sprintf("%d/%s", itemID, objectName),
		src,
		-1,
		minio.PutObjectOptions{})
	if err != nil {
		return &info, objectName, err
	}

	log.Printf("%#v\n", info)

	return &info, objectName, err
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
