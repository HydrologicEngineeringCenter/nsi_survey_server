package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/stores"
	"github.com/jackc/pgx"

	"github.com/labstack/echo/v4"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
)

type SurveyHandler struct {
	store         *stores.SurveyStore
	surveyEventId int
}

func CreateSurveyHandler(ss *stores.SurveyStore, eventId int) *SurveyHandler {
	sh := SurveyHandler{
		store:         ss,
		surveyEventId: eventId,
	}
	return &sh
}

func (sh *SurveyHandler) GetSurveyReport(c echo.Context) error {
	var eventId int
	eventId, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		eventId = sh.surveyEventId
	}
	s, err := sh.store.GetReport(eventId)
	if err != nil {
		return err
	}
	headers := "srId, userId, userName,completed,isControl,saId,fdId,x,y,invalidStructure,noStreetView,cbfips,occtype,stDamcat,foundHt,numStory,sqft,foundType,rsmeansType,quality,constType,garage,roofStyle\r\n"

	resp := c.Response()
	resp.Header().Set("Content-type", "text/csv")
	resp.Header().Set("Content-Disposition", "attachment; filename=surveys.csv")
	resp.Header().Set("Pragma", "no-cache")
	resp.Header().Set("Expires", "0")
	w := resp.Writer
	w.Write([]byte(headers))
	for _, record := range s {
		vals := record.String()
		for i, val := range vals {
			if i > 0 {
				w.Write([]byte(","))
			}
			if _, err := w.Write([]byte(val)); err != nil {
				log.Println("error writing headers to csv:", err)
				return err
			}
		}
		w.Write([]byte("\r\n"))
	}
	return err
}

func (sh *SurveyHandler) GetSurvey(c echo.Context) error {
	claims := c.Get("NSIUSER").(models.JwtClaim)
	userId := claims.Sub
	assignmentInfo, err := sh.store.GetAssignmentInfo(userId, sh.surveyEventId)
	if err != nil {
		return err
	}
	structure := models.SurveyStructure{}
	if assignmentInfo.Completed == nil || *assignmentInfo.Completed { //anything other than an explicit 'false'
		nextSurvey := assignmentInfo.NextSurvey

		if assignmentInfo.NextControl != nil && *assignmentInfo.NextControl < *assignmentInfo.NextSurvey {
			nextSurvey = assignmentInfo.NextControl
		}
		saId, err := sh.store.AssignSurvey(userId, *nextSurvey)
		if err != nil {
			log.Printf("Error assigning Survey: %s", err)
			pgerr := err.(pgx.PgError)
			if pgerr.Code == "23503" && pgerr.TableName == "survey_assignment" {
				return c.String(200, `{"result":"completed"}`) //this should only occur when we are out of surveys
			}
			return err
		}

		structure, err = sh.store.GetStructure(*nextSurvey, saId)
		if err != nil {
			return err
		}

	} else {
		structure, err = sh.store.GetStructure(*assignmentInfo.SE_ID, *assignmentInfo.SA_ID)
		if err != nil {
			return err
		}
	}
	return c.JSON(http.StatusOK, structure)
}

func (sh *SurveyHandler) SaveSurvey(c echo.Context) error {
	s := models.SurveyStructure{}
	if err := c.Bind(&s); err != nil {
		return err
	}
	err := sh.store.SaveSurvey(&s)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, `{"result":"success"}`)
}
