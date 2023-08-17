package raddoc

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/request"
	"cpm-rad-backend/domain/response"
	"cpm-rad-backend/domain/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getByItemFunc func(context.Context, SearchSpec, uint) (Response, int64, error)

func (fn getByItemFunc) GetByItem(ctx context.Context, spec SearchSpec, ID uint) (Response, int64, error) {
	return fn(ctx, spec, ID)
}

func GetByItemHandler(svc getByItemFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		search := ParseSearchSpec(c)
		data, total, err := svc.GetByItem(c.Request().Context(), search, utils.StringToUint(c.Param("itemid")))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, response.DataDoc[Item, listOfDoc]{
			Item:  data.Item,
			DOC:   data.ListOfDoc,
			Page:  search.GetPage(),
			Limit: search.GetLimit(),
			Total: total,
		})
	}
}

func ParseSearchSpec(c echo.Context) SearchSpec {
	return SearchSpec{
		Pagination: request.GetPagination(
			utils.StringToInt(c.QueryParam("page")),
			utils.StringToInt(c.QueryParam("limit")),
		),
		SeqNo:          c.QueryParam("searchRowNo"),
		InvNo:          c.QueryParam("searchInvNo"),
		Qty:            c.QueryParam("searchQuantity"),
		Arrival:        c.QueryParam("searchArrival"),
		Inspection:     c.QueryParam("searchInspection"),
		CreateBy:       c.QueryParam("searchCreater"),
		StateName:      c.QueryParam("searchStateName"),
		SortSeqNo:      c.QueryParam("sortRowNo"),
		SortInvNo:      c.QueryParam("sortInvNo"),
		SortQty:        c.QueryParam("sortQuantity"),
		SortArrival:    c.QueryParam("sortArrival"),
		SortInspection: c.QueryParam("sortInspection"),
		SortCreateBy:   c.QueryParam("sortCreater"),
		SortStateName:  c.QueryParam("sortStateName"),
	}
}
