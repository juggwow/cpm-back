package form

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/response"
	"fmt"
	"net/http"

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
