package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/db"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("M4_SECRET")
	jwtMiddleware := middleware.JWT([]byte(jwtSecret))

	// Middleware to ensure the user is still marked as logged in in DB.
	checkLoggedIn := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user")
			if user == nil {
				return c.JSON(http.StatusUnauthorized, "Missing or invalid token")
			}
			token := user.(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			username, ok := claims["name"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "Invalid token claims")
			}
			dbUser, err := db.FindOneUser(username)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "Invalid user")
			}
			if !dbUser.IsLoggedIn {
				return c.JSON(http.StatusUnauthorized, "User logged out")
			}
			return next(c)
		}
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", routes.Login)

	// Logout route - ensure token is valid and user still logged in
	e.POST("/logout", routes.Logout, jwtMiddleware, checkLoggedIn)

	// Register route
	e.POST("/register", routes.Register)

	// Get user notes
	r := e.Group("/notes")
	r.Use(jwtMiddleware)
	r.Use(checkLoggedIn)
	r.GET("/mynotes", routes.MyNotes)
	r.POST("/addnote", routes.AddNote)

	// Healthcheck route
	e.GET("/healthcheck", routes.Healthcheck)

	e.Logger.Fatal(e.Start(":9051"))

	defer func() {
		cerr := db.Disconnect()
		if err == nil {
			err = cerr
		}
	}()
}
