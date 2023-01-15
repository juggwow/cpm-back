package employee

import (
	"context"
	"cpm-rad-backend/domain/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getByIDFunc func(context.Context, string) (Employee, error)

func (fn getByIDFunc) GetByID(ctx context.Context, employeeID string) (Employee, error) {
	return fn(ctx, employeeID)
}

// GetByIDHandler godoc
// @Summary get an employee by id
// @Param employeeId path string true "Employee ID"
// @Tags Employees
// @Accept json
// @Produce json
// @Success 200 {object} EmployeeResponse
// @Router /api/v1/employees/{employeeId} [get]
// @Security ApiKeyAuth
func GetByIDHandler(svc getByIDFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		employeeID := c.Param("employeeId")
		log := logger.Unwrap(c)

		employee, err := svc.GetByID(c.Request().Context(), employeeID)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, employee.ToResponse())
	}
}
