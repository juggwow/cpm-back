package health_check

import (
	"cpm-rad-backend/domain/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Summary Check server health.
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {string} ok
// @Router /healths [get]
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, config.DBPass)
}
