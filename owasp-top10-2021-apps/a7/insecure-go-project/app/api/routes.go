package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// HealthCheck is the heath check function.
func HealthCheck(c echo.Context) error {
	// test123
	return c.String(http.StatusOK, "WORKING!\n")
}
