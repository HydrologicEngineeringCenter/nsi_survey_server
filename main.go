package main

import (
	"log"

	"github.com/USACE/microauth"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/auth"
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/config"
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/handlers"
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/stores"
)

// Config holds all runtime configuration provided via environment variables

const urlPrefix = "nsisapi"

func main() {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	cfg.SkipJWT = true

	ss, err := stores.CreateSurveyStore(&cfg)
	if err != nil {
		log.Printf("Unable to connect to database during startup: %s", err)
	}

	surveyHandler := handlers.CreateSurveyHandler(ss)
	jwtAuth := microauth.Auth{
		AuthMiddleware: auth.Appauth,
	}
	jwtAuth.LoadVerificationKey(cfg.Ippk)

	e := echo.New()

	e.Use(jwtAuth.AuthorizeMiddleware)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public Routes
	e.GET(urlPrefix+"/version", surveyHandler.Version)
	e.GET(urlPrefix+"/survey", surveyHandler.GetSurvey)
	e.POST(urlPrefix+"/survey", surveyHandler.CreateNewSurvey)
	//e.POST(urlPrefix+"/survey/:surveyId", surveyHandler.)
	//e.GET(urlPrefix+"/reports/surveys/:eventID", surveyHandler.GetSurveyReport)
	//e.GET(urlPrefix+"/survey/create", surveyHandler.CreateNewSurvey)

	//new endpoints
	// - create new survey
	// - add list of survey elements
	// - assign users to survey

	e.Logger.Fatal(e.Start(":" + cfg.Port))

}
