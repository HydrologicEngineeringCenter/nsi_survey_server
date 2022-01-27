package main

import (
	"log"

	"github.com/USACE/microauth"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	. "github.com/HydrologicEngineeringCenter/nsi_survey_server/auth"
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
	auth := microauth.Auth{
		AuthRoute: Appauth,
		Aud:       cfg.Aud,
		Store:     ss,
	}
	auth.LoadVerificationKey(cfg.Ippk)

	e := echo.New()

	//e.Use(jwtAuth.AuthorizeMiddleware)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET(urlPrefix+"/version", surveyHandler.Version)
	e.GET(urlPrefix+"/surveys", auth.AuthorizeRoute(surveyHandler.GetSurveysForUser, PUBLIC))
	e.POST(urlPrefix+"/survey", auth.AuthorizeRoute(surveyHandler.CreateNewSurvey, ADMIN, PUBLIC))
	e.PUT(urlPrefix+"/survey/:surveyid", auth.AuthorizeRoute(surveyHandler.UpdateSurvey, ADMIN, SURVEY_OWNER))
	e.GET(urlPrefix+"/survey/:surveyid/members", auth.AuthorizeRoute(surveyHandler.GetSurveyMembers, ADMIN, SURVEY_OWNER))
	e.POST(urlPrefix+"/survey/:surveyid/member", auth.AuthorizeRoute(surveyHandler.UpsertSurveyMember, ADMIN, SURVEY_OWNER))
	e.DELETE(urlPrefix+"/survey/member/:memberid", auth.AuthorizeRoute(surveyHandler.RemoveSurveyMember, ADMIN, SURVEY_OWNER))
	e.DELETE(urlPrefix+"/survey/:surveyid/member/:memberid", auth.AuthorizeRoute(surveyHandler.RemoveMemberFromSurvey, ADMIN, SURVEY_OWNER))
	e.POST(urlPrefix+"/survey/:surveyid/elements", auth.AuthorizeRoute(surveyHandler.InsertSurveyElements, ADMIN, SURVEY_OWNER))
	e.POST(urlPrefix+"/survey/:surveyid/assignments", auth.AuthorizeRoute(surveyHandler.AddAssignments, ADMIN))
	e.GET(urlPrefix+"/survey/:surveyid/assignment", auth.AuthorizeRoute(surveyHandler.AssignSurveyElement, SURVEY_MEMBER))
	e.POST(urlPrefix+"/survey/assignment", auth.AuthorizeRoute(surveyHandler.SaveSurveyAssignment, SURVEY_MEMBER))
	e.GET(urlPrefix+"/users/search", auth.AuthorizeRoute(surveyHandler.SearchUsers, PUBLIC))
	e.GET(urlPrefix+"/survey/:surveyid/report", auth.AuthorizeRoute(surveyHandler.GetSurveyReport, ADMIN, SURVEY_OWNER))

	e.Logger.Fatal(e.Start(":" + cfg.Port))

}
