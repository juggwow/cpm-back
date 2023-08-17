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

		search := ParseProgressReportSearch(c)
		data, total, err := svc.GetProgressReport(c.Request().Context(), search, utils.StringToUint(c.Param("id")))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, response.Data[ProgressReport]{
			Data:  data,
			Page:  search.GetPage(),
			Limit: search.GetLimit(),
			Total: total,
		})
	}
}

func ParseProgressReportSearch(c echo.Context) ProgressReportSearch {
	return ProgressReportSearch{

		Pagination: request.GetPagination(
			utils.StringToInt(c.QueryParam("page")),
			utils.StringToInt(c.QueryParam("limit")),
		),
		SequencesNo:     c.QueryParam("searchRowNo"),
		ItemName:        c.QueryParam("searchItemName"),
		Invoice:         c.QueryParam("searchInvNo"),
		Arrival:         c.QueryParam("searchArrival"),
		Inspection:      c.QueryParam("searchInspection"),
		StateName:       c.QueryParam("searchStateName"),
		SortSequencesNo: c.QueryParam("sortRowNo"),
		SortItemName:    c.QueryParam("sortItemName"),
		SortInvoice:     c.QueryParam("sortInvNo"),
		SortArrival:     c.QueryParam("sortArrival"),
		SortInspection:  c.QueryParam("sortInspection"),
		SortStateName:   c.QueryParam("sortStateName"),
	}
}

type getCheckReportFunc func(context.Context, CheckReportSearch, uint) (CheckReports, int64, error)

func (fn getCheckReportFunc) GetCheckReport(ctx context.Context, search CheckReportSearch, ID uint) (CheckReports, int64, error) {
	return fn(ctx, search, ID)
}

func GetCheckReportHandler(svc getCheckReportFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		search := ParseCheckSearch(c)
		data, total, err := svc.GetCheckReport(c.Request().Context(), search, utils.StringToUint(c.Param("id")))
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, response.Error{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, response.Data[CheckReport]{
			Data:  data,
			Page:  search.GetPage(),
			Limit: search.GetLimit(),
			Total: total,
		})
	}
}

func ParseCheckSearch(c echo.Context) CheckReportSearch {

	// var employee *auth.JwtEmployeeClaims
	// if claims, err := auth.GetAuthorizedClaims(c); err == nil {
	// 	employee = &claims
	// }

	return CheckReportSearch{

		Pagination: request.GetPagination(
			utils.StringToInt(c.QueryParam("page")),
			utils.StringToInt(c.QueryParam("limit")),
		),
		SequencesNo:     c.QueryParam("searchRowNo"),
		ItemName:        c.QueryParam("searchItemName"),
		Invoice:         c.QueryParam("searchInvNo"),
		Arrival:         c.QueryParam("searchArrival"),
		Inspection:      c.QueryParam("searchInspection"),
		Amount:          c.QueryParam("searchAmount"),
		Good:            c.QueryParam("searchGood"),
		Waste:           c.QueryParam("searchWaste"),
		SortSequencesNo: c.QueryParam("sortRowNo"),
		SortItemName:    c.QueryParam("sortItemName"),
		SortInvoice:     c.QueryParam("sortInvNo"),
		SortArrival:     c.QueryParam("sortArrival"),
		SortInspection:  c.QueryParam("sortInspection"),
		SortAmount:      c.QueryParam("sortAmount"),
		SortGood:        c.QueryParam("sortGood"),
		SortWaste:       c.QueryParam("sortWaste"),
	}

}
