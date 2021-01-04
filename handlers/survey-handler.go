package handlers

import (
	"fmt"
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
	fmt.Println(assignmentInfo)
	structure := models.SurveyStructure{}
	if assignmentInfo.Completed == nil || *assignmentInfo.Completed { //anything other than 'false'
		nextSurvey := assignmentInfo.NextSurvey
		if assignmentInfo.NextControl < assignmentInfo.NextSurvey {
			nextSurvey = assignmentInfo.NextControl
		}
		saId, err := sh.store.AssignSurvey(userId, nextSurvey)
		if err != nil {
			return err
		}

		structure, err = sh.store.GetStructure(nextSurvey, saId)
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
