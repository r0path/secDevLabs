package routes

import (
	"net/http"
	"errors"

	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m2/cool_games/server/app/db"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m2/cool_games/server/app/types"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m2/cool_games/server/app/util"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

// Register tries to register a new user into the database
func Register(c echo.Context) error {
	u := new(types.RequestUser)
	if err := c.Bind(u); err != nil {
		return err
	}

	attemptUsername := u.Username
	attemptPassword := u.Password

	_, err := db.FindOneUser(attemptUsername)
	if err == nil {
		return c.JSON(http.StatusConflict, "User already exists!")
	}

	// If the error is something other than "no documents found", it's a real DB error
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(http.StatusInternalServerError, "Database error while checking username, please try again later")
		}
	}

	hashedPassword, salt, err := util.Hash(attemptPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error hashing user password, please try again later")
	}

	newUser := types.User{attemptUsername, hashedPassword, salt}

	err = db.InsertOneUser(newUser)
	if err != nil {
		// handle duplicate key just in case a race created the user after the previous check
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, we := range writeErr.WriteErrors {
				if we.Code == 11000 {
					return c.JSON(http.StatusConflict, "User already exists!")
				}
			}
		}
		return c.JSON(http.StatusInternalServerError, "Error registering new user, please try again later")
	}

	return c.JSON(http.StatusOK, "User successfully registered!")
}
