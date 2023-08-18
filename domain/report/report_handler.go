package report

import (
	"context"
	"cpm-rad-backend/domain/auth"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/utils"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type createFunc func(context.Context, RequestReportCreate, File) (ResponseReportDetail, error)

func (fn createFunc) Create(ctx context.Context, r RequestReportCreate, f File) (ResponseReportDetail, error) {
	return fn(ctx, r, f)
}

func CreateHandler(svc createFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		claims, _ := auth.GetAuthorizedClaims(c)
		log.Info(strings.Join([]string{claims.EmployeeID, claims.FirstName, claims.LastName}, " "))

		// check permission

		var r RequestReportCreate
		r.ItemID = c.FormValue("itemID")
		r.Arrival = c.FormValue("arrival")
		r.Inspection = c.FormValue("inspection")
		r.TaskMaster = c.FormValue("taskMaster")
		r.Invoice = c.FormValue("invoice")
		r.Quantity = c.FormValue("quantity")
		r.Country = c.FormValue("country")
		r.Brand = c.FormValue("brand")
		r.Model = c.FormValue("model")
		r.Serial = c.FormValue("serial")
		r.PeaNo = c.FormValue("peano")
		r.CreateBy = strings.Join([]string{claims.FirstName, claims.LastName}, " ") //c.FormValue("createby")
		r.Status = c.FormValue("status")

		if err := invalidRequest(&r); err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, utils.ReaponseError{Error: err.Error()})
		}

		form, _ := c.MultipartForm()
		files := form.File["filesAttach"]

		if err := invalidFile(files); err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, utils.ReaponseError{Error: err.Error()})
		}

		fileType := func() []string {
			if data := c.FormValue("docType"); data != "" {
				return strings.Split(data, ",")
			}
			return []string{}
		}

		res, err := svc.Create(c.Request().Context(), r, File{
			Info: files,
			Type: fileType(),
		})

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, utils.ReaponseError{Error: err.Error()})
		}
		return c.JSON(http.StatusCreated, res)
	}
}

type getFunc func(context.Context, uint) (ResponseReportDetail, error)

func (fn getFunc) Get(ctx context.Context, id uint) (ResponseReportDetail, error) {
	return fn(ctx, id)
}

func GetHandler(svc getFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Error(err.Error())
			return c.String(http.StatusBadRequest, fmt.Sprintf("require id : %s", err.Error()))
		}

		res, err := svc.Get(c.Request().Context(), uint(id))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, utils.ReaponseError{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, res)
	}
}

type updateFunc func(context.Context, RequestReportUpdate, File) (ResponseReportDetail, error)

func (fn updateFunc) Update(ctx context.Context, r RequestReportUpdate, f File) (ResponseReportDetail, error) {
	return fn(ctx, r, f)
}

func UpdateHandler(svc updateFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		log := logger.Unwrap(c)

		claims, _ := auth.GetAuthorizedClaims(c)
		log.Info(strings.Join([]string{claims.EmployeeID, claims.FirstName, claims.LastName}, " "))

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, utils.ReaponseError{Error: err.Error()})
		}

		var r RequestReportUpdate
		r.ID = uint(id)
		r.ItemID = c.FormValue("itemID")
		r.Arrival = c.FormValue("arrival")
		r.Inspection = c.FormValue("inspection")
		r.TaskMaster = c.FormValue("taskMaster")
		r.Invoice = c.FormValue("invoice")
		r.Quantity = c.FormValue("quantity")
		r.Country = c.FormValue("country")
		r.Brand = c.FormValue("brand")
		r.Model = c.FormValue("model")
		r.Serial = c.FormValue("serial")
		r.PeaNo = c.FormValue("peano")
		r.Status = c.FormValue("status")
		r.UpdateBy = strings.Join([]string{claims.FirstName, claims.LastName}, " ")

		//// varidate data before update

		// if utils.IsEmpty(r.ItemID) {
		// 	return errors.New(":Invalid Data Request")
		// }

		// if err := invalidRequest(&r); err != nil {
		// 	log.Error(err.Error())
		// 	return c.JSON(http.StatusBadRequest, utils.ReaponseError{Error: err.Error()})
		// }

		form, _ := c.MultipartForm()
		files := form.File["filesAttach"]

		if err := invalidFile(files); err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, utils.ReaponseError{Error: err.Error()})
		}

		var updateFile []UpdateFile
		json.Unmarshal([]byte(c.FormValue("changeFileType")), &updateFile)
		// fmt.Printf("update : %+v", updateFile)

		fileType := func() []string {
			if data := c.FormValue("docType"); data != "" {
				return strings.Split(data, ",")
			}
			return []string{}
		}

		fileDel := func() []string {
			if data := c.FormValue("removeFile"); data != "" {
				return strings.Split(data, ",")
			}
			return []string{}
		}
		fmt.Println(len(fileType()))
		fmt.Println(len(fileDel()))
		res, err := svc.Update(c.Request().Context(), r, File{
			Info:   files,
			Type:   fileType(),
			Update: updateFile,
			Delete: fileDel(),
		})

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, utils.ReaponseError{Error: err.Error()})
		}

		// claims, _ := auth.GetAuthorizedClaims(c)
		// log.Info(claims.EmployeeID)

		// jobID, err := svc.Create(c.Request().Context(), reqJob, claims.EmployeeID)

		// err := svc.Update(c.Request().Context(), req)
		// if err != nil {
		// 	log.Error(err.Error())
		// 	return c.JSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		// }

		return c.JSON(http.StatusCreated, res)
	}
}

func invalidRequest(r *RequestReportCreate) error {
	if utils.IsEmpty(r.ItemID) {
		return errors.New(":Invalid Data Request")
	}

	return nil
}

func invalidFile(files []*multipart.FileHeader) error {
	const LIMIT int64 = 52428800
	for _, file := range files {
		if (file.Header["Content-Type"][0] != "application/pdf") && (LIMIT < file.Size) {
			return errors.New(":Invalid File Request")
		}
	}
	return nil
}
