package main

import (
	"log"
	"os"

	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m2/cool_games/server/app/db"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m2/cool_games/server/app/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// JWT middleware: enforce authentication on all routes except the public ones
	jwtSecret := os.Getenv("M4_SECRET")
	jwtConfig := middleware.JWTConfig{
		SigningKey:  []byte(jwtSecret),
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		Skipper: func(c echo.Context) bool {
			p := c.Path()
			return p == "/login" || p == "/register" || p == "/healthcheck"
		},
	}
	e.Use(middleware.JWTWithConfig(jwtConfig))

	// Login route
	e.POST("/login", routes.Login)

	// Register route
	e.POST("/register", routes.Register)

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
