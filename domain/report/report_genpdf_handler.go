package report

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type genPdfFunc func(context.Context, uint) (FileResponse, error)

func (fn genPdfFunc) GenPdf(ctx context.Context, id uint) (FileResponse, error) {
	return fn(ctx, id)
}

func GenPdfHandler(svc genPdfFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		reportID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Error(err.Error())
			return c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		}

		res, err := svc.GenPdf(c.Request().Context(), uint(reportID))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		}

		// c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", res.Name))
		return c.Blob(http.StatusOK, res.Ext, res.Obj)
	}
}
