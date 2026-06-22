package routes

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/db"
	"github.com/labstack/echo"
)

// MyNotes returns all user's notes.
func MyNotes(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)

	// Verify that the user's session is still active. Logout updates the IsLoggedIn
	// flag in the database, but previously issued JWTs remain valid until expiration.
	// Check the flag to ensure tokens are effectively revoked on logout.
	u, err := db.FindOneUser(username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid user")
	}

	if !u.IsLoggedIn {
		return c.JSON(http.StatusUnauthorized, "Token revoked or user logged out")
	}

	notes := db.FindNotes(username)

	return c.JSON(http.StatusOK, notes)
}
