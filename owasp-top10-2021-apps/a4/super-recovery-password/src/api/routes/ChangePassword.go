package routes

import (
	"api/database"
	"api/services"
	"api/types"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func ChangePassword(c echo.Context) (err error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*services.JwtCustomClaims)
	if claims.Recovery != true {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "invalid token",
		})
	}

	// Bind and validate incoming password payload
	u := new(types.ChangePassword)
	if err = c.Bind(u); err != nil {
		return
	}
	if u.Password != u.RepeatPassword {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "password don`t match",
		})
	}

	// Verify that the presented recovery token matches the server-side stored token for this user.
	// The token string is available in jwt.Token.Raw when parsed by middleware.
	tokenString := user.Raw
	valid, err := database.CheckRecoveryToken(claims.Name, tokenString)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to validate token",
		})
	}
	if !valid {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "invalid token",
		})
	}

	password := types.ChangePassword{
		Password:       u.Password,
		RepeatPassword: u.RepeatPassword,
	}

	// Change the password
	err = database.ChangePassword(claims.Name, password.Password, password.RepeatPassword)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusOK, echo.Map{
			"message": "failed to change password",
		})
	}

	// Invalidate the recovery token so it cannot be reused. This is best-effort; log on failure but do not fail the request
	if err := database.ClearRecoveryToken(claims.Name); err != nil {
		fmt.Println("failed to clear recovery token:", err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
