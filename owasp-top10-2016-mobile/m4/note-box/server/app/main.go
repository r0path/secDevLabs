package main

import (
	"log"
	"os"

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

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", routes.Login)

	// Logout route
	e.POST("/logout", routes.Logout, jwtMiddleware)

	// Register route
	e.POST("/register", routes.Register)

	// Get user notes
	r := e.Group("/notes")
	r.Use(jwtMiddleware)
	r.GET("/mynotes", routes.MyNotes)
	r.POST("/addnote", routes.AddNote)

	// Healthcheck route
	e.GET("/healthcheck", routes.Healthcheck)

	// Start server with TLS if cert and key are provided via environment variables.
	cert := os.Getenv("M4_TLS_CERT")
	key := os.Getenv("M4_TLS_KEY")
	if cert != "" && key != "" {
		log.Println("Starting server with TLS on :9051")
		e.Logger.Fatal(e.StartTLS(":9051", cert, key))
	}
	log.Println("TLS certificates not provided; starting server without TLS (insecure). Set M4_TLS_CERT and M4_TLS_KEY to enable TLS.")
	e.Logger.Fatal(e.Start(":9051"))

	defer func() {
		cerr := db.Disconnect()
		if err == nil {
			err = cerr
		}
	}()
}
