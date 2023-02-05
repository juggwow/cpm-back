package form

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/response"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type createFunc func(context.Context, Request, string) (uint, error)

func (fn createFunc) Create(ctx context.Context, req Request, createdBy string) (uint, error) {
	return fn(ctx, req, createdBy)
}

func CreateHandler(svc createFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req Request
		log := logger.Unwrap(c)
		if err := c.Bind(&req); err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		}

		if invalidRequest(&req) {
			return c.JSON(http.StatusBadRequest, response.Error{Error: fmt.Sprint(req)})
		}

		// claims, _ := auth.GetAuthorizedClaims(c)
		// jobID, err := svc.Create(c.Request().Context(), reqJob, claims.EmployeeID)
		formID, err := svc.Create(c.Request().Context(), req, req.CreateBy)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		}

		return c.JSON(http.StatusCreated, &response.ID{ID: formID})
	}
}

func invalidRequest(req *Request) bool {
	// if req.ItemID == 0 {
	// 	return true
	// }

	return req.ItemID == 0
}

type getFunc func(context.Context, uint) (Response, error)

func (fn getFunc) Get(ctx context.Context, id uint) (Response, error) {
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
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, res)
	}
}

type getCountryFunc func(context.Context) (Countrys, error)

func (fn getCountryFunc) GetCountry(ctx context.Context) (Countrys, error) {
	return fn(ctx)
}

func GetCountryHandler(svc getCountryFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		countrys, err := svc.GetCountry(c.Request().Context())
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, countrys)
	}
}

type fileUploadFunc func(context.Context, *multipart.FileHeader, int) (FileUploadResponse, error)

func (fn fileUploadFunc) FileUpload(ctx context.Context, file *multipart.FileHeader, itemID int) (FileUploadResponse, error) {
	return fn(ctx, file, itemID)
}

func FileUploadHandler(svc fileUploadFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var fus FileUploadResponses
		log := logger.Unwrap(c)

		itemID, err := strconv.Atoi(c.Param("itemid"))
		if err != nil {
			log.Error(err.Error())
			return c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		}
		fieldName := c.Param("fieldName")

		// file, err := c.FormFile(fieldName)

		form, _ := c.MultipartForm()
		files := form.File[fieldName]
		// filePaths := []string{}
		for _, file := range files {
			// check file type pdf and size < 50 MB
			src, err := file.Open()
			if err != nil {
				return c.JSON(http.StatusBadRequest, response.Error{Error: err.Error()})
			}
			defer src.Close()
			res, err := svc.FileUpload(c.Request().Context(), file, itemID)
			fus = append(fus, res)
			if err != nil {
				log.Error(err.Error())
				return c.JSON(http.StatusBadRequest, response.Error{Error: err.Error()})
			}
		}

		return c.JSON(http.StatusOK, fus)
	}
}

type getDocTypeFunc func(context.Context) (DocTypes, error)

func (fn getDocTypeFunc) GetDocType(ctx context.Context) (DocTypes, error) {
	return fn(ctx)
}

func GetDocTypeHandler(svc getDocTypeFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		docTypes, err := svc.GetDocType(c.Request().Context())
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, docTypes)
	}
}

type fileDeleteFunc func(context.Context, int, string) error

func (fn fileDeleteFunc) FileDelete(ctx context.Context, itemID int, objectName string) error {
	return fn(ctx, itemID, objectName)
}

func FileDeleteHandler(svc fileDeleteFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// var fus FileUploadResponses
		log := logger.Unwrap(c)

		itemID, err := strconv.Atoi(c.Param("itemid"))
		if err != nil {
			log.Error(err.Error())
			return c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		}
		objectName := c.QueryParam("obj")

		err = svc.FileDelete(c.Request().Context(), itemID, objectName)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		}

		return c.String(http.StatusOK, "success")
	}
}
