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

var newSurveyId string = "b4dd9464-a0a6-4a9e-b793-935920dccc46"

func TestCreateSurvey(t *testing.T) {
	createJSON := `{"title":"Survey Test","description":"This is a description of the test survey","active":true}`
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
	updateJSON := fmt.Sprintf(`{"id":"%s","title":"Survey Test Updated","description":"This is a description of survey edited","active":false}`, newSurveyId)
	rec, c := buildContext(http.MethodPost, updateJSON)
	h := buildHandler(t)
	if assert.NoError(t, h.UpdateSurvey(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestAddSurveyOwner(t *testing.T) {
	t.Log(newSurveyId)
	payload := fmt.Sprintf(`{"surveyId":"%s","userId":"987654"}`, newSurveyId)
	rec, c := buildContext(http.MethodPost, payload)
	h := buildHandler(t)
	if assert.NoError(t, h.AddSurveyOwner(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestAddSecondSurveyOwner(t *testing.T) {
	t.Log(newSurveyId)
	payload := fmt.Sprintf(`{"surveyId":"%s","userId":"887654"}`, newSurveyId)
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
	payload := fmt.Sprintf(`
	[
		{"surveyId":"%s","surveyOrder":1,"fdId":95009, "isControl":false},
		{"surveyId":"%s","surveyOrder":2,"fdId":95008, "isControl":false},
		{"surveyId":"%s","surveyOrder":3,"fdId":95007, "isControl":false},
		{"surveyId":"%s","surveyOrder":4,"fdId":95006, "isControl":false},
		{"surveyId":"%s","surveyOrder":5,"fdId":95005, "isControl":true},
		{"surveyId":"%s","surveyOrder":6,"fdId":95004, "isControl":false},
		{"surveyId":"%s","surveyOrder":7,"fdId":95003, "isControl":true},
		{"surveyId":"%s","surveyOrder":8,"fdId":95002, "isControl":false},
		{"surveyId":"%s","surveyOrder":9,"fdId":95001, "isControl":false},
	]`, newSurveyId, newSurveyId, newSurveyId, newSurveyId, newSurveyId, newSurveyId, newSurveyId, newSurveyId, newSurveyId)
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

func TestGetSurveyAssignment(t *testing.T) {
	t.Log(newSurveyId)
	rec, c := buildContext(http.MethodGet, "")
	c.SetParamNames("surveyID")
	c.SetParamValues(newSurveyId)
	h := buildHandler(t)
	if assert.NoError(t, h.GetSurvey(c)) {
		fmt.Println(rec.Body.String())
		assert.Equal(t, http.StatusOK, rec.Code)
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
		Sub:      "987654",
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
