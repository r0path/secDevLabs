package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/globocom/secDevLabs/owasp-top10-2021-apps/a1/ecommerce-api/app/db"
	"github.com/labstack/echo"
)

// HealthCheck is the heath check function.
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING\n")
}

// GetTicket returns the userID ticket.
func GetTicket(c echo.Context) error {
	id := c.Param("id")

	cookie, err := c.Cookie("sessionIDa5")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"result": "error", "details": "Unauthorized."})
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"result": "error", "details": "Unauthorized."})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"result": "error", "details": "Unauthorized."})
	}

	username, _ := claims["name"].(string)
	if strings.TrimSpace(username) == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"result": "error", "details": "Unauthorized."})
	}

	userDataQuery := map[string]interface{}{"userID": id, "username": username}
	userDataResult, err := db.GetUserData(userDataQuery)
	if err != nil {
		// could not find this user in MongoDB (or MongoDB err connection)
		return c.JSON(http.StatusNotFound, map[string]string{"result": "error", "details": "Ticket not found."})
	}

	format := c.QueryParam("format")
	if format == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"result":   "success",
			"username": userDataResult.Username,
			"ticket":   userDataResult.Ticket,
		})
	}

	msgTicket := fmt.Sprintf("Hey, %s! This is your ticket: %s\n", userDataResult.Username, userDataResult.Ticket)
	return c.String(http.StatusOK, msgTicket)
}
