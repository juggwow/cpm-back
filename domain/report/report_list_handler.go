package report

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/request"
	"cpm-rad-backend/domain/response"
	"cpm-rad-backend/domain/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getProgressReportFunc func(context.Context, ProgressReportSearch, uint) (ProgressReports, int64, error)

func (fn getProgressReportFunc) GetProgressReport(ctx context.Context, search ProgressReportSearch, ID uint) (ProgressReports, int64, error) {
	return fn(ctx, search, ID)
}

func GetProgressReportHandler(svc getProgressReportFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		search := ParseSearch(c)
		data, total, err := svc.GetProgressReport(c.Request().Context(), search, utils.StringToUint(c.Param("id")))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}

		// data := progressReport.ToResponse()

		return c.JSON(http.StatusOK, response.Data[ProgressReport]{
			Data:  data,
			Page:  search.GetPage(),
			Limit: search.GetLimit(),
			Total: total,
		})
	}
}

func ParseSearch(c echo.Context) ProgressReportSearch {

	// var employee *auth.JwtEmployeeClaims
	// if claims, err := auth.GetAuthorizedClaims(c); err == nil {
	// 	employee = &claims
	// }

	return ProgressReportSearch{
		Pagination:      request.GetPagination(utils.StringToInt(c.QueryParam("page")), utils.StringToInt(c.QueryParam("limit"))),
		SequencesNo:     c.QueryParam("seq"),
		ItemName:        c.QueryParam("itemName"),
		Invoice:         c.QueryParam("invNo"),
		Arrival:         c.QueryParam("arrival"),
		Inspection:      c.QueryParam("inspection"),
		StateName:       c.QueryParam("stateName"),
		SortSequencesNo: c.QueryParam("sSeq"),
		SortItemName:    c.QueryParam("sItemName"),
		SortInvoice:     c.QueryParam("sInvNo"),
		SortArrival:     c.QueryParam("sArrival"),
		SortInspection:  c.QueryParam("sInspection"),
		SortStateName:   c.QueryParam("sStateName"),
	}
}
