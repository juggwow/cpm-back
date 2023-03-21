package report

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/utils"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type createFunc func(context.Context, Report, string) (uint, error)

func (fn createFunc) Create(ctx context.Context, r Report, createdBy string) (uint, error) {
	return fn(ctx, r, createdBy)
}

func CreateHandler(svc createFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)
		// var req Request
		// log := logger.Unwrap(c)
		// if err := c.Bind(&req); err != nil {
		var r Report
		bind(c, &r)

		if err := invalidRequest(&r); err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, utils.ReaponseError{Error: err.Error()})
		}

		docTypes := strings.Split(c.FormValue("docType"), ",")

		form, _ := c.MultipartForm()
		files := form.File["filesAttach"]
		// filePaths := []string{}
		var f AttachFile
		for i, file := range files {
			// check file type pdf and size < 50 MB
			// src, err := file.Open()
			// if err != nil {
			// 	return c.JSON(http.StatusBadRequest, response.Error{Error: err.Error()})
			// }
			// defer src.Close()
			// res, err := svc.FileUpload(c.Request().Context(), file, itemID)
			// fus = append(fus, res)
			// if err != nil {
			// 	log.Error(err.Error())
			// 	return c.JSON(http.StatusBadRequest, response.Error{Error: err.Error()})
			// }
			f.Name = file.Filename
			f.Size = strconv.FormatInt(file.Size, 10)
			f.Path = file.Header["Content-Type"][0]
			f.Type = utils.StringToUint(docTypes[i])
			r.AttachFiles = append(r.AttachFiles, f)
		}

		// var req Request
		// log := logger.Unwrap(c)
		// if err := c.Bind(&req); err != nil {
		// 	log.Error(err.Error())
		// 	return c.JSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		// }

		// if invalidRequest(&req) {
		// 	return c.JSON(http.StatusBadRequest, response.Error{Error: fmt.Sprint(req)})
		// }

		// // claims, _ := auth.GetAuthorizedClaims(c)
		// // jobID, err := svc.Create(c.Request().Context(), reqJob, claims.EmployeeID)
		// formID, err := svc.Create(c.Request().Context(), req, req.CreateBy)
		// if err != nil {
		// 	log.Error(err.Error())
		// 	return c.JSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		// }

		return c.JSON(http.StatusCreated, r)
	}
}

func bind(c echo.Context, r *Report) {

	r.ItemID = c.FormValue("itemID")
	r.Arrival = c.FormValue("arrival")
	r.Inspection = c.FormValue("inspection")
	r.TaskMaster = c.FormValue("taskMaster")
	r.Invoice = c.FormValue("invoice")
	r.Quantity = c.FormValue("quantity")
	r.Country = c.FormValue("country")
	r.Manufacturer = c.FormValue("manufacturer")
	r.Model = c.FormValue("model")
	r.Serial = c.FormValue("serial")
	r.PeaNo = c.FormValue("peano")
	r.CreateBy = c.FormValue("createby")
	r.Status = c.FormValue("status")

}

func invalidRequest(r *Report) error {
	if utils.IsEmpty(r.ItemID) {
		return errors.New(":Invalid Data Request")
	}

	return nil
}
