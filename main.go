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

	e.GET(urlPrefix+"/version", surveyHandler.Version)
	e.GET(urlPrefix+"/surveys", surveyHandler.GetSurveysForUser)
	e.POST(urlPrefix+"/survey", surveyHandler.CreateNewSurvey)
	e.PUT(urlPrefix+"/survey", surveyHandler.UpdateSurvey)
	e.GET(urlPrefix+"/survey/:surveyid/members", surveyHandler.GetSurveyMembers)
	e.POST(urlPrefix+"/survey/member", surveyHandler.UpsertSurveyMember)
	e.DELETE(urlPrefix+"/survey/member/:memberid", surveyHandler.RemoveSurveyMember)
	e.POST(urlPrefix+"/survey/elements", surveyHandler.InsertSurveyElements)
	e.POST(urlPrefix+"/survey/assignments", surveyHandler.AddAssignments)
	e.GET(urlPrefix+"/survey/:surveyid/assignment", surveyHandler.AssignSurveyElement)
	e.POST(urlPrefix+"/survey/assignment", surveyHandler.SaveSurveyAssignment)
	e.GET(urlPrefix+"/users/search", surveyHandler.SearchUsers)
	e.GET(urlPrefix+"/survey/:surveyid/report", surveyHandler.GetSurveyReport)

	e.Logger.Fatal(e.Start(":" + cfg.Port))

}
