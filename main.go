package main

import (
	"log"
	"net/http"

	"github.com/apex/gateway"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/auth"
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

	surveyHandler := handlers.CreateSurveyHandler(ss, cfg.SurveyEvent)
	jwtAuth := auth.Auth{
		Store: ss,
	}
	jwtAuth.LoadVerificationKey(cfg.Ippk)

	e := echo.New()

	e.Use(jwtAuth.Authorize)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public Routes
	e.GET("nsisapi/survey", surveyHandler.GetSurvey)
	e.POST("nsisapi/survey", surveyHandler.SaveSurvey)
	e.GET("nsisapi/reports/surveys/:eventID", surveyHandler.GetSurveyReport)

	if cfg.LambdaContext {
		log.Print("starting server; Running On AWS LAMBDA")
		log.Fatal(gateway.ListenAndServe("localhost:3030", e))
	} else {
		log.Print("starting server on port 3031")
		log.Fatal(http.ListenAndServe("0.0.0.0:3031", e))
	}
}
