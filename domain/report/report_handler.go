package report

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/utils"
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

		var r RequestReportCreate
		bind(c, &r)

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

		res, err := svc.Create(c.Request().Context(), r, File{
			Info: files,
			Type: strings.Split(c.FormValue("docType"), ","),
		})

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, utils.ReaponseError{Error: err.Error()})
		}
		return c.JSON(http.StatusCreated, res)
	}
}

func bind(c echo.Context, r *RequestReportCreate) {

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
	r.CreateBy = c.FormValue("createby")
	r.Status = c.FormValue("status")

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

// type updateFunc func(context.Context, Report, File) (Report, error)

// func (fn updateFunc) Update(ctx context.Context, r Report, f File) (Report, error) {
// 	return fn(ctx, r, f)
// }

// func UpdateHandler(svc updateFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		log := logger.Unwrap(c)

// 		id, err := strconv.Atoi(c.Param("id"))
// 		if err != nil {
// 			log.Error(err.Error())
// 			return c.JSON(http.StatusBadRequest, utils.ReaponseError{Error: err.Error()})
// 		}

// 		var r Report
// 		r.ID = uint(id)
// 		bind(c, &r)

// 		if err := invalidRequest(&r); err != nil {
// 			log.Error(err.Error())
// 			return c.JSON(http.StatusBadRequest, utils.ReaponseError{Error: err.Error()})
// 		}

// 		form, _ := c.MultipartForm()
// 		files := form.File["filesAttach"]

// 		if err := invalidFile(files); err != nil {
// 			log.Error(err.Error())
// 			return c.JSON(http.StatusBadRequest, utils.ReaponseError{Error: err.Error()})
// 		}

// 		var updateFile []UpdateFile
// 		json.Unmarshal([]byte(c.FormValue("updateDocType")), &updateFile)
// 		fmt.Printf("update : %+v", updateFile)

// 		r, err = svc.Update(c.Request().Context(), r, File{
// 			Info:   files,
// 			Type:   strings.Split(c.FormValue("docType"), ","),
// 			Update: updateFile,
// 			Delete: strings.Split(c.FormValue("delFile"), ","),
// 		})

// 		if err != nil {
// 			log.Error(err.Error())
// 			return c.JSON(http.StatusInternalServerError, utils.ReaponseError{Error: err.Error()})
// 		}
// 		return c.JSON(http.StatusCreated, r)

// 		// claims, _ := auth.GetAuthorizedClaims(c)
// 		// jobID, err := svc.Create(c.Request().Context(), reqJob, claims.EmployeeID)

// 		// err := svc.Update(c.Request().Context(), req)
// 		// if err != nil {
// 		// 	log.Error(err.Error())
// 		// 	return c.JSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
// 		// }

// 		// return c.String(http.StatusOK, "success")
// 	}
// }

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
