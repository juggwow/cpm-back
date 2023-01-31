package raddoc

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/request"
	"cpm-rad-backend/domain/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type getByItemFunc func(context.Context, SearchSpec, uint) (Response, int64, error)

func (fn getByItemFunc) GetByItem(ctx context.Context, spec SearchSpec, ID uint) (Response, int64, error) {
	return fn(ctx, spec, ID)
}

func GetByItemHandler(svc getByItemFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		ID, _ := strconv.Atoi(c.Param("itemid"))
		spec := ParseSearchSpec(c)
		data, total, err := svc.GetByItem(c.Request().Context(), spec, uint(ID))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}

		// data := res.ToResponse()

		return c.JSON(http.StatusOK, response.DataDoc[Item, listOfDoc]{
			Item:  data.Item,
			DOC:   data.ListOfDoc,
			Page:  spec.GetPage(),
			Limit: spec.GetLimit(),
			Total: total,
		})
	}
}

func ParseSearchSpec(c echo.Context) SearchSpec {

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// var employee *auth.JwtEmployeeClaims
	// if claims, err := auth.GetAuthorizedClaims(c); err == nil {
	// 	employee = &claims
	// }

	seqNo, _ := strconv.Atoi(c.QueryParam("seqNo"))
	return SearchSpec{
		Pagination:     request.GetPagination(page, limit),
		SeqNo:          seqNo,
		InvNo:          "",
		Qty:            "",
		Arrival:        "",
		Inspection:     "",
		CreateBy:       "",
		StateName:      "",
		SortSeqNo:      "",
		SortInvNo:      "",
		SortQty:        "",
		SortArrival:    "",
		SortInspection: "",
		SortCreateBy:   "",
		SortStateName:  "",
	}
}
