package report

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/utils"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type createFunc func(context.Context, Report, File) (Report, error)

func (fn createFunc) Create(ctx context.Context, r Report, f File) (Report, error) {
	return fn(ctx, r, f)
}

func CreateHandler(svc createFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		var r Report
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

		r, err := svc.Create(c.Request().Context(), r, File{
			Info: files,
			Type: strings.Split(c.FormValue("docType"), ","),
		})

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, utils.ReaponseError{Error: err.Error()})
		}
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

func invalidFile(files []*multipart.FileHeader) error {
	const LIMIT int64 = 52428800
	for _, file := range files {
		if (file.Header["Content-Type"][0] != "application/pdf") && (LIMIT < file.Size) {
			return errors.New(":Invalid File Request")
		}
	}
	return nil
}
