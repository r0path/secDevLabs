package main

import (
	"fmt"
	"os"
	"strconv"
	"net/http"
	"time"

	"github.com/globocom/secDevLabs/owasp-top10-2021-apps/a7/insecure-go-project/app/api"
	"github.com/globocom/secDevLabs/owasp-top10-2021-apps/a7/insecure-go-project/app/config"
	db "github.com/globocom/secDevLabs/owasp-top10-2021-apps/a7/insecure-go-project/app/db/mongo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func main() {

	fmt.Println("[*] Starting Insecure Go Project...")

	// loading viper
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		errorAPI(err)
	}
	if err := viper.Unmarshal(&config.APIconfiguration); err != nil {
		errorAPI(err)
	}

	// check if MongoDB is accessible and credentials received are working.
	if _, err := checkMongoDB(); err != nil {
		fmt.Println("[X] ERROR MONGODB: ", err)
		os.Exit(1)
	}

	fmt.Println("[*] MongoDB: OK!")
	fmt.Println("[*] Viper loaded: OK!")

	echoInstance := echo.New()
	echoInstance.HideBanner = true

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())

	echoInstance.GET("/healthcheck", api.HealthCheck)
	APIport := fmt.Sprintf(":%d", getAPIPort())

	// Set reasonable server timeouts to mitigate slowloris-style DoS attacks.
	// ReadTimeout: maximum duration for reading the entire request, including the body.
	// WriteTimeout: maximum duration before timing out writes of the response.
	// IdleTimeout: maximum amount of time to wait for the next request when keep-alives are enabled.
	srv := &http.Server{
		Addr:         APIport,
		Handler:      echoInstance,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	echoInstance.Logger.Fatal(srv.ListenAndServe())
}

func errorAPI(err error) {
	fmt.Println("[x] Error starting Insecure Go Project:")
	fmt.Println("[x]", err)
	os.Exit(1)
}

func getAPIPort() int {
	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		apiPort = 10002
	}
	return apiPort
}

func checkMongoDB() (*db.DB, error) {
	return db.Connect()
}
