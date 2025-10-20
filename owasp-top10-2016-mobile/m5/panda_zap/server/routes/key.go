package routes

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type keyHolderV1 struct {
	Value int `bson:"key" json:"key"`
}

type keyHolderV2 struct {
	Value string `bson:"key" json:"key"`
}

// GetKeyV1 returns the encryption key
func (es *EchoServer) GetKeyV1(c echo.Context) error {
	// Access to the global encryption key is forbidden.
	// Returning 403 to prevent disclosure of sensitive information.
	return c.JSON(http.StatusForbidden, map[string]string{"error": "access to key is forbidden"})
}
