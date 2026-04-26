package routes

import (
	"crypto/rand"
	"encoding/binary"
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

	if es.MessageKey.Value == 0 {
		var b [2]byte
		if _, err := rand.Read(b[:]); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate key"})
		}
		es.MessageKey.Value = int(binary.BigEndian.Uint16(b[:])%100) + 1
	}

	return c.JSON(http.StatusOK,
		map[string]string{"key": strconv.Itoa(es.MessageKey.Value)})
}
