package routes

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/db"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/types"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/util"
	"github.com/labstack/echo"
)

// Register tries to register a new user into the database
func Register(c echo.Context) error {
	u := new(types.RequestUser)
	if err := c.Bind(u); err != nil {
		return err
	}

	attemptUsername := strings.TrimSpace(u.Username)
	attemptPassword := strings.TrimSpace(u.Password)

	// Basic server-side validation
	if attemptUsername == "" || attemptPassword == "" {
		return c.JSON(http.StatusBadRequest, "Username and password must not be empty")
	}

	// Username: 3-30 chars, letters, numbers, dot, underscore, hyphen
	var usernameRegex = regexp.MustCompile(`^[A-Za-z0-9._-]{3,30}$`)
	if !usernameRegex.MatchString(attemptUsername) {
		return c.JSON(http.StatusBadRequest, "Username must be 3-30 characters and contain only letters, numbers, dot, underscore or hyphen")
	}

	// Password: minimum length 8
	if len(attemptPassword) < 8 {
		return c.JSON(http.StatusBadRequest, "Password must be at least 8 characters long")
	}

	_, err := db.FindOneUser(attemptUsername)
	if err == nil {
		return c.JSON(http.StatusConflict, "User already exists!")
	}

	hashedPassword, salt, err := util.Hash(attemptPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error hashing user password, please try again later")
	}

	newUser := types.User{attemptUsername, hashedPassword, salt, false}

	err = db.InsertOneUser(newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error registering new user, please try again later")
	}

	return c.JSON(http.StatusOK, "User successfully registered!")
}
