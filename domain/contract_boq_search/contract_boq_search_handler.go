package contract_boq_search

import (
	"context"
	"cpm-rad-backend/domain/contract_boq"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/request"
	"cpm-rad-backend/domain/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type getFunc func(context.Context, BoQSearchSpec, uint) (contract_boq.BoQItems, int64, error)

func (fn getFunc) Get(ctx context.Context, spec BoQSearchSpec, ID uint) (contract_boq.BoQItems, int64, error) {
	return fn(ctx, spec, ID)
}

// GetHandler godoc
// @Summary get boq items
// @Tags BOQ
// @Accept json
// @Produce json
// @Param _ query BoQSearchSpec false "BoQ items Search Query"
// @Success 200 {object} response.ResponseData[contract_boq.BoQResponse]
// @Router api/v1/contract/:id/boq [get]
// @Security ApiKeyAuth
func GetHandler(svc getFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		ID, err := strconv.Atoi(c.Param("id"))
		spec := ParseSearchSpec(c)
		boqItems, total, err := svc.Get(c.Request().Context(), spec, uint(ID))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.ResponseError{Error: err.Error()})
		}

		data := boqItems.ToResponse()

		return c.JSON(http.StatusOK, response.ResponseData[contract_boq.BoQResponse]{
			Data:  data,
			Page:  spec.GetPage(),
			Limit: spec.GetLimit(),
			Total: total,
		})
	}
}

func ParseSearchSpec(c echo.Context) BoQSearchSpec {

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// var employee *auth.JwtEmployeeClaims
	// if claims, err := auth.GetAuthorizedClaims(c); err == nil {
	// 	employee = &claims
	// }

	seqNo, _ := strconv.Atoi(c.QueryParam("seqNo"))
	return BoQSearchSpec{
		Pagination:       request.GetPagination(page, limit),
		SequencesNo:      seqNo,
		ItemNo:           c.QueryParam("num"),
		ItemName:         c.QueryParam("name"),
		ItemGroup:        c.QueryParam("group"),
		ItemQuantity:     c.QueryParam("quantity"),
		ItemDelivery:     c.QueryParam("delivery"),
		ItemReceive:      c.QueryParam("receive"),
		ItemDamage:       c.QueryParam("damage"),
		SortSequencesNo:  "",
		SortItemNo:       "",
		SortItemName:     "",
		SortItemGroup:    "",
		SortItemQuantity: "",
		SortItemDelivery: "",
		SortItemReceive:  "",
		SortItemDamage:   "",
	}
}
