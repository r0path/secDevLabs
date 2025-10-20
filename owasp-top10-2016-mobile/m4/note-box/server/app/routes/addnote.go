package routes

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/db"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/types"
	"github.com/labstack/echo"
)

// AddNote attempts to add a new note in the database
func AddNote(c echo.Context) error {
	userVal := c.Get("user")
	if userVal == nil {
		return c.JSON(http.StatusUnauthorized, "Missing user token")
	}
	user, ok := userVal.(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Invalid user token")
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Invalid token claims")
	}

	receivedNote := new(types.Note)
	if err := c.Bind(receivedNote); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error unmarshaling new note")
	}

	nameClaim, ok := claims["name"]
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Token missing name claim")
	}
	jwtUsername, ok := nameClaim.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Invalid name claim in token")
	}

	// Normalize usernames to avoid case-sensitivity mismatches
	normalizedJWT := strings.ToLower(strings.TrimSpace(jwtUsername))
	normalizedOwner := strings.ToLower(strings.TrimSpace(receivedNote.OwnerUsername))

	if normalizedJWT != normalizedOwner {
		return c.JSON(http.StatusForbidden, "You can't add notes as "+receivedNote.OwnerUsername+"! You're logged in as "+jwtUsername)
	}

	err := db.InsertOneNote(receivedNote)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error inserting new note into the database")
	}

	return c.JSON(http.StatusOK, "Note added successfully")
}
