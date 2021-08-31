package stores

import (
	"context"
	"log"
	"strings"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/config"
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
	"github.com/google/uuid"
	"github.com/usace/goquery"
)

var NoResults string = "no rows in result set"

type SurveyStore struct {
	DS goquery.DataStore
}

func CreateSurveyStore(appConfig *config.Config) (*SurveyStore, error) {
	dbconf := appConfig.Rdbmsconfig()
	ds, err := goquery.NewRdbmsDataStore(&dbconf)
	if err != nil {
		log.Printf("Unable to connect to database during startup: %s", err)
	}
	log.Printf("Connected as %s to database %s:%s/%s", appConfig.Dbuser, appConfig.Dbhost, appConfig.Dbport, appConfig.Dbname)

	//ds.SetMaxOpenConns(4)
	ss := SurveyStore{ds}
	return &ss, nil
}

func (ss *SurveyStore) GetSurvey(surveyId uuid.UUID) (models.Survey, error) {
	survey := models.Survey{}
	err := ss.DS.Select().
		DataSet(&surveyTable).
		StatementKey("selectById").
		Dest(&survey).
		Params(surveyId).
		Fetch()
	return survey, err
}

func (ss *SurveyStore) CreateNewSurvey(survey models.Survey, userId string) (uuid.UUID, error) {
	var surveyId uuid.UUID
	err := goquery.Transaction(ss.DS, func(tx goquery.Tx) {
		err := ss.DS.Select().
			DataSet(&surveyTable).
			Tx(&tx).
			StatementKey("insert").
			Params(survey.Title, survey.Description, survey.Active).
			Dest(&surveyId).
			Fetch()

		if err != nil {
			panic(err)
		}
		ptx := tx.PgxTx()
		_, err = ptx.Exec(context.Background(), surveyTable.Statements["insert-owner"], surveyId, userId)
		if err != nil {
			panic(err)
		}
	})
	return surveyId, err
}

func (ss *SurveyStore) UpdateSurvey(survey models.Survey) error {
	err := ss.DS.Exec(goquery.NoTx, surveyTable.Statements["update"], survey.Title, survey.Description, survey.Active, survey.ID)
	return err
}

func (ss *SurveyStore) AddSurveyOwner(owner models.SurveyOwner) error {
	err := ss.DS.Exec(goquery.NoTx, surveyOwnerTable.Statements["insert"], owner.SurveyID, owner.UserID)
	return err
}

func (ss *SurveyStore) RemoveSurveyOwner(id uuid.UUID) error {
	err := ss.DS.Exec(goquery.NoTx, surveyOwnerTable.Statements["remove"], id)
	return err
}

func (ss SurveyStore) InsertSurveyElements(elements *[]models.SurveyElement) error {
	err := ss.DS.Insert(&surveyElementTable).
		Records(elements).
		Execute()

	if err != nil {
		log.Printf("Error inserting survey elements: %s", err)
	}
	return err
}

func (ss *SurveyStore) AssignSurvey(userId string, seId uuid.UUID) (uuid.UUID, error) {
	var saId uuid.UUID
	err := ss.DS.Select(surveyAssignmentTable.Statements["assignSurvey"]).
		Params(seId, userId).
		Dest(&saId).
		Fetch()
	return saId, err
}

func (ss SurveyStore) InsertSurveyAssignments(assignments *[]models.SurveyAssignment) error {
	err := ss.DS.Insert(&surveyAssignmentTable).
		Records(assignments).
		Execute()

	if err != nil {
		log.Printf("Error inserting survey assignments: %s", err)
	}
	return err
}

func (ss *SurveyStore) GetReport(surveyId uuid.UUID) ([]models.SurveyResult, error) {
	s := []models.SurveyResult{}
	err := ss.DS.Select(miscQueries.Statements["surveyReport"]).
		Params(surveyId).
		Dest(s).
		Fetch()
	return s, err
}

func (ss *SurveyStore) GetAssignmentInfo(userId string, surveyId uuid.UUID) (models.AssignmentInfo, error) {
	ai := models.AssignmentInfo{}
	err := ss.DS.Select(surveyAssignmentTable.Statements["assignmentInfo"]).
		Params(userId, surveyId).
		Dest(&ai).
		Fetch()
	if err != nil {
		return models.AssignmentInfo{}, err
	}
	return ai, err
}

func (ss *SurveyStore) GetFirstSurveyInEvent(surveyId uuid.UUID) (uuid.UUID, error) {
	var firstSurvey uuid.UUID
	err := ss.DS.Select("select id from survey_element where survey_order=(select min(survey_order) from survey_element where survey_event_id=$1)").
		Params(surveyId).
		Dest(firstSurvey).
		Fetch()
	return firstSurvey, err
}

func (ss *SurveyStore) GetStructure(seId uuid.UUID, saId uuid.UUID) (models.SurveyStructure, error) {
	s := models.SurveyStructure{}
	err := ss.DS.Select(surveyTable.Statements["survey"]).
		Params(saId).
		Dest(&s).
		Fetch()
	if err != nil {
		if err.Error() == NoResults {
			//no existing survey result, get survey data from nsi
			err := ss.DS.Select(surveyTable.Statements["nsi-survey"]).
				Params(seId, saId).
				Dest(&s).
				Fetch()
			if err != nil {
				log.Printf("Failed to retrieve structure: %s/n", err)
				return s, err
			}
			log.Printf("Returning NSI Data for survey assignment: %d/n", saId)
			s.OccupancyType = strings.Split(s.OccupancyType, "-")[0]
			return s, nil
		} else {
			log.Printf("Failed to query survey results for existing assignment: %s/n", err)
			return s, err //return error
		}
	}
	log.Printf("Returning existing Survey Result for survey assignment: %d/n", saId)
	return s, err //return survey from survey_result
}

/*
func (ss *SurveyStore) SaveSurvey(survey *models.SurveyStructure) error {
	err := goquery.Transaction(ss.DS, func(tx goquery.Tx) {
		_, txerr := tx.NamedExec(tables.Statements["upsertSurveyStructure"], survey)
		if txerr != nil {
			panic(txerr)
		}
		tx.MustExec(tables.Statements["updateAssignment"], survey.SAID)
	})
	return err
}
*/
