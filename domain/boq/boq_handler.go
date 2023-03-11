package boq

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/request"
	"cpm-rad-backend/domain/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type getFunc func(context.Context, ItemSearchSpec, uint) (Items, int64, error)

func (fn getFunc) Get(ctx context.Context, spec ItemSearchSpec, ID uint) (Items, int64, error) {
	return fn(ctx, spec, ID)
}

func GetHandler(svc getFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		ID, _ := strconv.Atoi(c.Param("id"))
		spec := ParseSearchSpec(c)
		boqItems, total, err := svc.Get(c.Request().Context(), spec, uint(ID))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}

		data := boqItems.ToResponse()

		return c.JSON(http.StatusOK, response.Data[Response]{
			Data:  data,
			Page:  spec.GetPage(),
			Limit: spec.GetLimit(),
			Total: total,
		})
	}
}

func ParseSearchSpec(c echo.Context) ItemSearchSpec {

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// var employee *auth.JwtEmployeeClaims
	// if claims, err := auth.GetAuthorizedClaims(c); err == nil {
	// 	employee = &claims
	// }

	// seqNo, _ := strconv.Atoi(c.QueryParam("seqNo"))
	return ItemSearchSpec{
		Pagination:       request.GetPagination(page, limit),
		SequencesNo:      c.QueryParam("order"),
		ItemNo:           c.QueryParam("itemNo"),
		ItemName:         c.QueryParam("name"),
		ItemGroup:        c.QueryParam("group"),
		ItemQuantity:     c.QueryParam("quantity"),
		ItemDelivery:     c.QueryParam("delivery"),
		ItemReceive:      c.QueryParam("good"),
		ItemDamage:       c.QueryParam("bad"),
		SortSequencesNo:  c.QueryParam("sorder"),
		SortItemNo:       c.QueryParam("sitemNo"),
		SortItemName:     c.QueryParam("sname"),
		SortItemGroup:    c.QueryParam("sgroup"),
		SortItemQuantity: c.QueryParam("squantity"),
		SortItemDelivery: c.QueryParam("sdelivery"),
		SortItemReceive:  c.QueryParam("sgood"),
		SortItemDamage:   c.QueryParam("sbad"),
	}
}

type getItemByIDFunc func(context.Context, uint) (ItemResponse, error)

func (fn getItemByIDFunc) GetItemByID(ctx context.Context, ID uint) (ItemResponse, error) {
	return fn(ctx, ID)
}

func GetItemByIDHandler(svc getItemByIDFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		ID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: logger.INVALID})
		}
		item, err := svc.GetItemByID(c.Request().Context(), uint(ID))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: logger.NOT_FOUND})
		}

		return c.JSON(http.StatusOK, item)
	}
}
