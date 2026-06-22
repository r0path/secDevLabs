package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/globocom/secDevLabs/owasp-top10-2021-apps/a3/copy-n-paste/app/types"
	"github.com/globocom/secDevLabs/owasp-top10-2021-apps/a3/copy-n-paste/app/util"
)

//HealthCheck is de health check function
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING!")
}

//Login is the login function
func Login(c echo.Context) error {

	loginAttempt := types.LoginAttempt{}
	err := c.Bind(&loginAttempt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "error", "details": "Error binding login attempt."})
	}

	validUser, err := util.AuthenticateUser(loginAttempt.User, loginAttempt.Pass)
	if err != nil {
		msgUser := err.Error()
		return c.JSON(http.StatusBadRequest, msgUser)
	}

	if validUser {
		msgUser := fmt.Sprintf("Welcome, %s!", loginAttempt.User)

		// create JWT token so the client can authenticate subsequent requests
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			// fallback secret for local/testing environments - should be set in production
			secret = "changeme"
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": loginAttempt.User,
			"exp":  time.Now().Add(24 * time.Hour).Unix(),
		})
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"result": "error", "details": "Error generating token."})
		}

		// set token as a secure, httpOnly cookie (works well for browser clients)
		cookie := &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		}
		c.SetCookie(cookie)

		// also return token in Authorization header for API clients
		c.Response().Header().Set("Authorization", "Bearer "+tokenString)

		return c.String(http.StatusOK, msgUser)
	}

	msgUser := fmt.Sprintf("User not found or wrong password!")
	return c.String(http.StatusBadRequest, msgUser)
}

//Register is the function to register a new user on bd
func Register(c echo.Context) error {

	RegisterAttempt := types.RegisterAttempt{}
	err := c.Bind(&RegisterAttempt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "error", "details": "Error binding register attempt."})
	}

	userCreated, err := util.NewUser(RegisterAttempt.User, RegisterAttempt.Pass, RegisterAttempt.PassCheck)
	if err != nil {
		msgUser := err.Error()
		return c.JSON(http.StatusOK, msgUser)
	}

	if userCreated {
		msgUser := fmt.Sprintf("User %s created!", RegisterAttempt.User)
		return c.String(http.StatusOK, msgUser)
	}

	msgUser := fmt.Sprintf("User already exists or passwords don't match!")
	return c.String(http.StatusOK, msgUser)
}
