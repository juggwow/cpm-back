package contract

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/response"
	"net/http"

	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type getByIDFunc func(context.Context, int) (Contract, error)

func (fn getByIDFunc) GetByID(ctx context.Context, ID int) (Contract, error) {
	return fn(ctx, ID)
}

func GetByIDHandler(svc getByIDFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		ID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		}

		contract, err := svc.GetByID(c.Request().Context(), ID)
		zap.L().Sugar().Infof("work contract %d", ID)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}
		data := contract.ToResponse()
		return c.JSON(http.StatusOK, data)
	}
}
