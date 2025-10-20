package routes

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m2/cool_games/server/app/db"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m2/cool_games/server/app/types"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m2/cool_games/server/app/util"
	"github.com/labstack/echo"
)

// Register tries to register a new user into the database
func Register(c echo.Context) error {
	u := new(types.RequestUser)
	if err := c.Bind(u); err != nil {
		return err
	}

	// Normalize and validate username/password
	username := strings.TrimSpace(u.Username)
	username = strings.ToLower(username)
	password := u.Password

	// Validate username length (3-30 characters)
	if username == "" || utf8.RuneCountInString(username) < 3 || utf8.RuneCountInString(username) > 30 {
		return c.JSON(http.StatusBadRequest, "Username must be between 3 and 30 characters and not empty")
	}

	// Validate password length (8-128 characters)
	if password == "" || utf8.RuneCountInString(password) < 8 || utf8.RuneCountInString(password) > 128 {
		return c.JSON(http.StatusBadRequest, "Password must be between 8 and 128 characters and not empty")
	}

	_, err := db.FindOneUser(username)
	if err == nil {
		return c.JSON(http.StatusConflict, "User already exists!")
	}

	hashedPassword, salt, err := util.Hash(password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error hashing user password, please try again later")
	}

	newUser := types.User{username, hashedPassword, salt}

	err = db.InsertOneUser(newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error registering new user, please try again later")
	}

	return c.JSON(http.StatusOK, "User successfully registered!")
}
