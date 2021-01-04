package main

import (
	"log"
	"net/http"

	"github.com/USACE/consequences-api/middleware"
	"github.com/apex/gateway"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/handlers"
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/stores"
)

// Config holds all runtime configuration provided via environment variables

func main() {
	var cfg models.Config
	if err := envconfig.Process("ns", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	cfg.SkipJWT = true

	ss, err := stores.CreateSurveyStore(&cfg)
	if err != nil {
		log.Printf("Unable to connect to database during startup: %s", err)
	}

	surveyHandler := handlers.CreateSurveyHandler(ss)

	e := echo.New()

	// Public Routes
	public := e.Group("")

	// Private Routes
	private := e.Group("")
	if cfg.SkipJWT == true {
		private.Use(middleware.MockIsLoggedIn)
	} else {
		private.Use(middleware.JWT, middleware.IsLoggedIn)
	}

	// Public Routes
	public.GET("nsi_api/survey_element", surveyHandler.GetSurvey)

	// Private Routes
	//private.POST("nsi_api/survey_result", handlers.PostSurveyResult)

	if cfg.LambdaContext {
		log.Print("starting server; Running On AWS LAMBDA")
		log.Fatal(gateway.ListenAndServe("localhost:3030", e))
	} else {
		log.Print("starting server on port 3031")
		log.Fatal(http.ListenAndServe("localhost:3031", e))
	}
}
