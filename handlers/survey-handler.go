package handlers

import (
	"net/http"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/stores"

	"github.com/labstack/echo/v4"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
)

type SurveyHandler struct {
	store *stores.SurveyStore
}

func CreateSurveyHandler(ss *stores.SurveyStore) *SurveyHandler {
	sh := SurveyHandler{
		store: ss,
	}
	return &sh
}

func (sh *SurveyHandler) GetSurvey(c echo.Context) error {
	userId := "rr"
	assignmentInfo, err := sh.store.GetAssignmentInfo(userId)
	if err != nil {
		return err
	}
	structure := models.NsiStructure{}
	if assignmentInfo.Completed == nil || *assignmentInfo.Completed { //anything other than 'false'
		nextSurvey := assignmentInfo.NextSurvey
		if assignmentInfo.NextControl < assignmentInfo.NextSurvey {
			nextSurvey = assignmentInfo.NextControl
		}
		structure, err = sh.store.GetStructure(nextSurvey)
		if err != nil {
			return err
		}
		err = sh.store.AssignSurvey(userId, nextSurvey)
		if err != nil {
			return err
		}
	} else {
		structure, err = sh.store.GetStructure(*assignmentInfo.SE_ID)
	}
	return c.JSON(http.StatusOK, structure)
}
