package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/stores"
	"github.com/USACE/microauth"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
)

var defaultUuid uuid.UUID

type SurveyHandler struct {
	store *stores.SurveyStore
}

func CreateSurveyHandler(ss *stores.SurveyStore) *SurveyHandler {
	sh := SurveyHandler{
		store: ss,
	}
	return &sh
}

func (sh *SurveyHandler) Version(c echo.Context) error {
	return c.String(http.StatusOK, "NSI Survey API Version 2.01 Development")
}

func (sh *SurveyHandler) CreateNewSurvey(c echo.Context) error {
	var survey = models.Survey{}
	if err := c.Bind(&survey); err != nil {
		return err
	}
	jwtclaims := c.Get("NSIUSER").(microauth.JwtClaim)

	newId, err := sh.store.CreateNewSurvey(survey, jwtclaims.Sub)
	if err != nil {
		log.Println("Error creating survey -----------")
		log.Println(err)
		log.Println(survey)
		log.Println("--------------------------------")
		return err
	}

	return c.JSONBlob(http.StatusCreated, []byte(fmt.Sprintf(`{"surveyId":"%s"}`, newId)))
}

func (sh *SurveyHandler) UpdateSurvey(c echo.Context) error {
	var survey = models.Survey{}
	if err := c.Bind(&survey); err != nil {
		return err
	}
	err := sh.store.UpdateSurvey(survey)
	if err != nil {
		log.Printf("Error updating survey: %s", err)
		return err
	}
	return c.String(http.StatusOK, "")
}

func (sh *SurveyHandler) GetSurveyMembers(c echo.Context) error {
	surveyId, err := uuid.Parse(c.Param("surveyid"))
	if err != nil {
		return err
	}
	members, err := sh.store.GetSurveyMembers(surveyId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &members)
}

func (sh *SurveyHandler) UpsertSurveyMember(c echo.Context) error {
	var surveyMember = models.SurveyMember{}
	if err := c.Bind(&surveyMember); err != nil {
		return err
	}
	err := sh.store.UpsertSurveyMember(surveyMember)
	if err != nil {
		log.Printf("Error adding survey member: %s", err)
		return err
	}
	return c.String(http.StatusCreated, "")
}

func (sh *SurveyHandler) RemoveSurveyMember(c echo.Context) error {
	memberId, err := uuid.Parse(c.Param("memberid"))
	if err != nil {
		return err
	}
	err = sh.store.RemoveSurveyMember(memberId)
	if err != nil {
		log.Printf("Error removing survey member: %s", err)
		return err
	}
	return c.String(http.StatusCreated, "")
}

func (sh *SurveyHandler) InsertSurveyElements(c echo.Context) error {
	var elements = []models.SurveyElement{}
	if err := c.Bind(&elements); err != nil {
		return err
	}
	err := sh.store.InsertSurveyElements(&elements)
	if err != nil {
		return err
	}
	return c.String(http.StatusCreated, "")
}

func (sh *SurveyHandler) AddAssignments(c echo.Context) error {
	var assignments = []models.SurveyAssignment{}
	if err := c.Bind(&assignments); err != nil {
		return err
	}
	err := sh.store.InsertSurveyAssignments(&assignments)
	if err != nil {
		return err
	}
	return c.String(http.StatusCreated, "")

}

func (sh *SurveyHandler) GetSurveyReport(c echo.Context) error {
	surveyId, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		return err
	}

	s, err := sh.store.GetReport(surveyId)
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

func (sh *SurveyHandler) GetSurveysForUser(c echo.Context) error {
	claims := c.Get("NSIUSER").(microauth.JwtClaim)
	userId := claims.Sub
	surveys, err := sh.store.GetSurveysforUser(userId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, surveys)
}

func (sh *SurveyHandler) GetSurvey(c echo.Context) error {
	surveyId, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		return err
	}
	claims := c.Get("NSIUSER").(microauth.JwtClaim)
	userId := claims.Sub
	assignmentInfo, err := sh.store.GetAssignmentInfo(userId, surveyId)
	if err != nil {
		return err
	}

	var structure models.SurveyStructure
	var nextSurvey *uuid.UUID
	if assignmentInfo.Completed == nil {
		//the user does not have any uncompleted surveys assigned.  get a new one.
		nextSurvey = assignmentInfo.NextSurveySEID
		if assignmentInfo.NextControlOrder != nil && assignmentInfo.NextSurveyOrder != nil &&
			*assignmentInfo.NextControlOrder < *assignmentInfo.NextSurveyOrder {
			nextSurvey = &assignmentInfo.NextControlSEID
		}
		if nextSurvey != nil {
			saId, err := sh.store.AssignSurvey(userId, *nextSurvey)
			fmt.Println(nextSurvey)
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
		}
	} else {
		structure, err = sh.store.GetStructure(*assignmentInfo.SEID, *assignmentInfo.SAID)
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

func (sh *SurveyHandler) SearchUsers(c echo.Context) error {

	q := c.QueryParam("q")
	r := c.QueryParam("r")
	p := c.QueryParam("p")

	rows, errRow := strconv.Atoi(r)
	page, errPage := strconv.Atoi(p)
	if q == "" || errRow != nil || errPage != nil {
		return errors.New("Invalid Query Parameters")
	}
	users, err := sh.store.DS.Select("select * from users where username like $1 limit $2 offset $3").
		Params("%"+q+"%", rows, rows*page).
		FetchJSON()
	if err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, users)
}
