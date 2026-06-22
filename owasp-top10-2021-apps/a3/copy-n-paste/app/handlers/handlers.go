package handlers

import (
	"fmt"
	"net/http"
	"time"

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
		// On successful authentication, set a simple session cookie so middleware can enforce protected routes.
		cookie := new(http.Cookie)
		cookie.Name = "session_user"
		cookie.Value = loginAttempt.User
		cookie.Path = "/"
		cookie.HttpOnly = true
		// For demo purposes we set a short expiration. In production, use secure (HTTPS-only), signed tokens.
		cookie.Expires = time.Now().Add(30 * time.Minute)
		c.SetCookie(cookie)

		msgUser := fmt.Sprintf("Welcome, %s!", loginAttempt.User)
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
