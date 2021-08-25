package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/stores"
	"github.com/USACE/microauth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/usace/goquery"
)

var newSurveyId string = "2c97b020-d185-453c-93b8-ab7879f4a620"

func TestCreateSurvey(t *testing.T) {
	createJSON := `{"title":"Survey 2","description":"This is a description of survey 2","active":true}`
	rec, c := buildContext(http.MethodPost, createJSON)
	h := buildHandler(t)
	if assert.NoError(t, h.CreateNewSurvey(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		out := rec.Body.String()
		t.Log(out)
		assert.Equal(t, `{"surveyId":`, out[0:12])
		newSurveyId = out[13 : len(out)-2]
		t.Log(newSurveyId)
	}
}

func TestUpdateSurvey(t *testing.T) {
	t.Log(newSurveyId)
	updateJSON := fmt.Sprintf(`{"id":"%s","title":"Survey Updated","description":"This is a description of survey edited","active":false}`, newSurveyId)
	rec, c := buildContext(http.MethodPost, updateJSON)
	h := buildHandler(t)
	if assert.NoError(t, h.UpdateSurvey(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestAddSurveyOwner(t *testing.T) {
	t.Log(newSurveyId)
	payload := fmt.Sprintf(`{"surveyId":"%s","userId":"9876543"}`, newSurveyId)
	rec, c := buildContext(http.MethodPost, payload)
	h := buildHandler(t)
	if assert.NoError(t, h.AddSurveyOwner(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestRemoveSurveyOwner(t *testing.T) {
	t.Log(newSurveyId)
	payload := ""
	rec, c := buildContext(http.MethodPost, payload)
	c.SetParamNames("surveyOwnerId")
	c.SetParamValues("fe311b7a-1b17-49a7-bf25-425240717c39")
	h := buildHandler(t)
	if assert.NoError(t, h.RemoveSurveyOwner(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestInsertSurveyElements(t *testing.T) {
	t.Log(newSurveyId)
	payload := `
	[
		{"surveyId":"2c97b020-d185-453c-93b8-ab7879f4a620","surveyOrder":1,"fdId":1234, "isControl":false},
		{"surveyId":"2c97b020-d185-453c-93b8-ab7879f4a620","surveyOrder":2,"fdId":1235, "isControl":true},
		{"surveyId":"2c97b020-d185-453c-93b8-ab7879f4a620","surveyOrder":3,"fdId":1236, "isControl":false}
	]`
	rec, c := buildContext(http.MethodPost, payload)
	h := buildHandler(t)
	if assert.NoError(t, h.InsertSurveyElements(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestInsertSurveyAssignments(t *testing.T) {
	t.Log(newSurveyId)
	payload := `
	[
		{"seId":"68900070-656f-4e94-abdd-8bc878eaa2dc","completed":false, "assignedTo":"987654"},
		{"seId":"68900070-656f-4e94-abdd-8bc878eaa2dc","completed":false, "assignedTo":"987655"}
	]`
	rec, c := buildContext(http.MethodPost, payload)
	h := buildHandler(t)
	if assert.NoError(t, h.AddAssignments(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

/////Private support methods///////

func getDataStore(t *testing.T) goquery.DataStore {
	config := goquery.RdbmsConfigFromEnv()
	ds, err := goquery.NewRdbmsDataStore(config)
	if err != nil {
		t.Errorf("Failed to connect to store:%s\n", err)
	}
	return ds
}

func getSurveyStore(t *testing.T) *stores.SurveyStore {
	datastore := getDataStore(t)
	return &stores.SurveyStore{datastore}
}

func getClaims() microauth.JwtClaim {
	claims := microauth.JwtClaim{
		Sub:      "12345678",
		Aud:      []string{"nsi-survey"},
		UserName: "Test User",
	}
	return claims
}

func buildHandler(t *testing.T) *SurveyHandler {
	surveystore := getSurveyStore(t)
	return CreateSurveyHandler(surveystore)
}

func buildContext(method string, payload string) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("NSIUSER", getClaims())
	return rec, c
}
