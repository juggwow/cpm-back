package boqItem

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/request"
	"cpm-rad-backend/domain/response"
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
		spec := ParseSearchSpec(c)
		data, total, err := svc.Get(c.Request().Context(), spec, uint(ID))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, response.Data[BoqItemList]{
			Data:  data,
			Page:  spec.GetPage(),
			Limit: spec.GetLimit(),
			Total: total,
		})
	}
}

func ParseSearchSpec(c echo.Context) SearchSpec {

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	return SearchSpec{
		Pagination:        request.GetPagination(page, limit),
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
