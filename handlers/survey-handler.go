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

const version = "2.0.1 Development"

type SurveyHandler struct {
	store *stores.SurveyStore
}

func CreateSurveyHandler(ss *stores.SurveyStore) *SurveyHandler {
	sh := SurveyHandler{
		store: ss,
	}
	return &sh
}

//Returns the API version as a text
//PUBLIC API
func (sh *SurveyHandler) Version(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("NSI Survey API Version %s", version))
}

//List the surveys that the requesting user (via the JWT Claim sub identifier) is a member of in a JSON array
//
//PUBLIC API
func (sh *SurveyHandler) GetSurveysForUser(c echo.Context) error {
	claims := c.Get("NSIUSER").(microauth.JwtClaim)
	userId := claims.Sub
	roles := claims.Roles
	var surveys *[]models.Survey
	var err error
	if microauth.Contains_string(roles, "ADMIN") {
		surveys, err = sh.store.GetSurveysforAdmin()
	} else {
		surveys, err = sh.store.GetSurveysforUser(userId)
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, surveys)
}

//Creates a new survey and returns the generated identifier in a JSON document
//
//e.g. {"surveyId":"1111-1111-111111"}
//
//PRIVATE API restricted to the ADMIN role
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

//Updates a survey and returns an empty HTTP OK result on success.
//
//PRIVATE API restricted to the ADMIN or SURVEY_OWNER roles
func (sh *SurveyHandler) UpdateSurvey(c echo.Context) error {
	var survey = models.Survey{}
	if err := c.Bind(&survey); err != nil {
		return err
	}
	if !validateUrl(survey.ID, c) {
		return errors.New("Invalid Request")
	}
	err := sh.store.UpdateSurvey(survey)
	if err != nil {
		log.Printf("Error updating survey: %s", err)
		return err
	}
	return c.String(http.StatusOK, "")
}

//Gets an array of survey members for a given survey. Returns a JSON array.
//
//PRIVATE API restricted to the ADMIN or SURVEY_OWNER roles
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

//Updates/Inserts a survey member record. Returns an empty HTTP CREATED (201) result on success.
//
//PRIVATE API restricted to the ADMIN or SURVEY_OWNER roles
func (sh *SurveyHandler) UpsertSurveyMember(c echo.Context) error {
	var surveyMember = models.SurveyMember{}
	if err := c.Bind(&surveyMember); err != nil {
		return err
	}
	if !validateUrl(surveyMember.SurveyID, c) {
		return errors.New("Invalid Request")
	}
	err := sh.store.UpsertSurveyMember(surveyMember)
	if err != nil {
		log.Printf("Error adding survey member: %s", err)
		return err
	}
	return c.String(http.StatusCreated, "")
}

//Removes a survey member record. Returns an empty HTTP OK result on success.
//
//PRIVATE API restricted to the ADMIN or SURVEY_OWNER roles
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
	return c.String(http.StatusOK, "")
}

//Removes a survey member record from a specific survey. Returns an empty HTTP OK result on success.
//
//PRIVATE API restricted to the ADMIN or SURVEY_OWNER roles
func (sh *SurveyHandler) RemoveMemberFromSurvey(c echo.Context) error {
	memberId := c.Param("memberid")
	surveyId, err := uuid.Parse(c.Param("surveyid"))
	if err != nil {
		return err
	}
	err = sh.store.RemoveMemberFromSurvey(memberId, surveyId)
	if err != nil {
		log.Printf("Error removing survey member: %s", err)
		return err
	}
	return c.String(http.StatusOK, "")
}

//Inserts an array of survey elements.  Returns an empty HTTP CREATED (201) result on success.
//
//PRIVATE API restricted to the ADMIN or SURVEY_OWNER roles
func (sh *SurveyHandler) InsertSurveyElements(c echo.Context) error {
	var elements = []models.SurveyElement{}
	if err := c.Bind(&elements); err != nil {
		return err
	}
	servId, ok := validateElements(&elements)
	if !ok || !validateUrl(servId, c) {
		return errors.New("Invalid Request")
	}

	err := sh.store.InsertSurveyElements(&elements)
	if err != nil {
		return err
	}
	return c.String(http.StatusCreated, "")
}

//method for manually making assignments to users.  Typically assignments should be made using the AssignSurveyElement method
//but this allows for admins to override the normal assignment algorithm. Returns an empty HTTP CREATED (201) result on success.
//
//PRIVATE API restricted to the ADMIN or SURVEY_OWNER roles
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

//Assigns a survey element to a survey member.  It works in the following manner:
//If a user has an existing assignment that has not been saved, then that survey is returned. If the user does not have an existing assignment,
//then surveys will be assigned in ascending order based on the survey order field.  Each survey will be
//assigned to a single user with the exception of control surveys.  Control surveys will be assigned to all users.
//When there are no more surveys to assign (all surveys are completed and the user has completed their control surveys),
//then the function will return an empty survey (e.g. id values of 0).
//
//PRIVATE API restricted to the SURVEY_MEMBER role
func (sh *SurveyHandler) AssignSurveyElement(c echo.Context) error {
	surveyId, err := uuid.Parse(c.Param("surveyid"))
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

//Saves the survey assignment and returns an HTTP OK on success
//
//PRIVATE API restricted to the SURVEY_MEMBER role

func (sh *SurveyHandler) SaveSurveyAssignment(c echo.Context) error {
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

//Search the user list.  This method takes three query parameters:
//
//q: the query term that will match against the user name
//
//r: the number of rows to return
//
//p: the page number to return
//
//PUBLIC API
func (sh *SurveyHandler) SearchUsers(c echo.Context) error {

	q := c.QueryParam("q")
	r := c.QueryParam("r")
	p := c.QueryParam("p")

	rows, errRow := strconv.Atoi(r)
	page, errPage := strconv.Atoi(p)
	if q == "" || errRow != nil || errPage != nil {
		return errors.New("Invalid Query Parameters")
	}
	users, err := sh.store.DS.Select("select * from users where user_name ilike $1 limit $2 offset $3").
		Params("%"+q+"%", rows, rows*page).
		FetchJSON()
	if err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, users)
}

//Validate that the survey name is available. Returns true if name is unused
//
//q: the query term that will match against the survey name
//
//PUBLIC API
func (sh *SurveyHandler) ValidSurveyName(c echo.Context) error {

	q := c.QueryParam("q")

	if q == "" {
		return errors.New("Invalid Query Parameters")
	}
	var surveys []models.Survey
	err := sh.store.DS.Select("select * from survey where title = $1").
		Params(q).
		Dest(&surveys).
		Fetch()
	if err != nil {
		return err
	}
	invalid := len(surveys) > 0
	return c.JSONBlob(http.StatusOK, []byte(`{"result":`+strconv.FormatBool(!invalid)+`}`))
}

//Returns a CSV dump of the survey results for a given survey
//
//PRIVATE API restructed to the ADMIN or SURVEY_OWNER role
func (sh *SurveyHandler) GetSurveyReport(c echo.Context) error {
	surveyId, err := uuid.Parse(c.Param("surveyid"))
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

func validateElements(elements *[]models.SurveyElement) (uuid.UUID, bool) {
	var surveyId uuid.UUID
	for i, v := range *elements {
		if i == 0 {
			surveyId = v.SurveyID
		} else {
			if surveyId != v.SurveyID {
				return surveyId, false
			}
		}
	}
	return surveyId, true
}

// Checks surveyId in body matches with surveyId passed by URI
// Do not use with handlers where surveyId isn't an expected URI param
func validateUrl(surveyId uuid.UUID, c echo.Context) bool {
	s := c.Get("NSISURVEY")
	if s == nil {
		return false
	}
	surveyUrlId, ok := s.(uuid.UUID)
	if !ok {
		return false
	}
	if surveyId == surveyUrlId {
		return true
	}
	log.Printf("Invalid Request.  URL SurveyId (%s) does not match Payload Survey Id(%s)", surveyUrlId, surveyId)
	return false
}
