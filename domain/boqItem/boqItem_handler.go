package boqItem

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/request"
	"cpm-rad-backend/domain/response"
	"cpm-rad-backend/domain/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type getFunc func(context.Context, SearchSpec, uint) (BoqItemLists, int64, error)

func (fn getFunc) Get(ctx context.Context, spec SearchSpec, ID uint) (BoqItemLists, int64, error) {
	return fn(ctx, spec, ID)
}

func GetHandler(svc getFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		ID, _ := strconv.Atoi(c.Param("id"))
		search := ParseSearchSpec(c)
		data, total, err := svc.Get(c.Request().Context(), search, uint(ID))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, response.Data[BoqItemList]{
			Data:  data,
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
		SearchRowNo:       c.QueryParam("searchRowNo"),
		SearchNumber:      c.QueryParam("searchNumber"),
		SearchGroupName:   c.QueryParam("searchGroupName"),
		SearchName:        c.QueryParam("searchName"),
		SearchQuantity:    c.QueryParam("searchQuantity"),
		SearchDeliveryQty: c.QueryParam("searchDeliveryQty"),
		SearchReceiveQty:  c.QueryParam("searchReceiveQty"),
		SearchDamageQty:   c.QueryParam("searchDamageQty"),
		SortRowNo:         c.QueryParam("sortRowNo"),
		SortNumber:        c.QueryParam("sortNumber"),
		SortGroupName:     c.QueryParam("sortGroupName"),
		SortName:          c.QueryParam("sortName"),
		SortQuantity:      c.QueryParam("sortQuantity"),
		SortDeliveryQty:   c.QueryParam("sortDeliveryQty"),
		SortReceiveQty:    c.QueryParam("sortReceiveQty"),
		SortDamageQty:     c.QueryParam("sortDamageQty"),
	}
}
