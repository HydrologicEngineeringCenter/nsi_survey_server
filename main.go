package main

import (
	"log"
	"net/http"

	"github.com/USACE/consequences-api/middleware"
	"github.com/apex/gateway"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/handlers"
)

// Config holds all runtime configuration provided via environment variables
type Config struct {
	SkipJWT       bool
	LambdaContext bool
	DBUser        string
	DBPass        string
	DBName        string
	DBHost        string
	DBSSLMode     string
}

func main() {
	var cfg Config
	if err := envconfig.Process("consequences", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	cfg.SkipJWT = true
	e := echo.New()

	// Public Routes
	//public := e.Group("")

	// Private Routes
	private := e.Group("")
	if cfg.SkipJWT == true {
		private.Use(middleware.MockIsLoggedIn)
	} else {
		private.Use(middleware.JWT, middleware.IsLoggedIn)
	}

	// Public Routes
	// NOTE: ALL GET REQUESTS ARE ALLOWED WITHOUT AUTHENTICATION USING JWTConfig Skipper. See appconfig/jwt.go

	// Next Element

	// Computes
	// public.GET("consequences/computes", handlers.ListComputes(db))
	// public.GET("consequences/computes/:compute_id", handlers.GetCompute(db))
	// public.GET("consequences/computes/:compute_id/result", handlers.GetComputeResult(db))
	private.GET("nsi_api/survey_element", handlers.GetNextElement())
	private.POST("nsi_api/survey_result", handlers.PostSurveyResult())

	if cfg.LambdaContext {
		log.Print("starting server; Running On AWS LAMBDA")
		log.Fatal(gateway.ListenAndServe("localhost:3030", e))
	} else {
		log.Print("starting server")
		log.Fatal(http.ListenAndServe("localhost:3031", e))
	}
}
