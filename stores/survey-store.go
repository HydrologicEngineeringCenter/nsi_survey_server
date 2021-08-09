package stores

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/usace/dataquery"
)

type SurveyStore struct {
	store dataquery.SqlDataStore
}

func CreateSurveyStore(appConfig *models.Config) (*SurveyStore, error) {
	dbconf := appConfig.Rdbmsconfig()
	con, err := dataquery.NewSqlConnection(&dbconf)
	if err != nil {
		log.Printf("Unable to connect to database during startup: %s", err)
	}
	log.Printf("Connected as %s to database %s:%s/%s", appConfig.Dbuser, appConfig.Dbhost, appConfig.Dbport, appConfig.Dbname)
	con.SetMaxOpenConns(4)
	ss := SurveyStore{
		store: dataquery.SqlDataStore{
			DB: con,
		},
	}
	return &ss, nil
}

func (ss *SurveyStore) GetAssignmentInfo(userId string, surveyEventId int) (models.AssignmentInfo, error) {
	ai := []models.AssignmentInfo{}
	err := ss.store.DB.Select(&ai, tables.Statements["assignmentInfo"], surveyEventId, userId, surveyEventId, surveyEventId, userId)
	if err != nil {
		return models.AssignmentInfo{}, err
	}
	if len(ai) == 0 {
		return models.AssignmentInfo{}, errors.New("Invalid Record")
	}
	if ai[0].NextSurvey == nil {
		ns, err := ss.GetFirstSurveyInEvent(surveyEventId)
		if err != nil {
			return models.AssignmentInfo{}, err
		}
		ai[0].NextSurvey = &ns
	}
	return ai[0], err
}

func (ss *SurveyStore) GetFirstSurveyInEvent(surveyEventId int) (int, error) {
	var firstSurvey int
	err := ss.store.DB.Get(&firstSurvey, "select min(id) from survey_element where survey_event_id=$1", surveyEventId)
	return firstSurvey, err
}

func (ss *SurveyStore) GetStructure(seId int, saId int) (models.SurveyStructure, error) {
	s := models.SurveyStructure{}
	err := ss.store.DB.Get(&s, tables.Statements["survey"], strconv.Itoa(saId))
	if err != nil {
		if err == sql.ErrNoRows {
			//no existing survey result, get survey data from nsi
			err := ss.store.DB.Get(&s, tables.Statements["nsi_survey"], seId, strconv.Itoa(saId))
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

func (ss *SurveyStore) AssignSurvey(userId string, seId int) (int, error) {
	var saId int
	err := ss.store.DB.QueryRow(tables.Statements["assignSurvey"], seId, userId).Scan(&saId)
	if err != nil {
		return -1, err
	}
	return int(saId), nil
}

func (ss *SurveyStore) SaveSurvey(survey *models.SurveyStructure) error {
	err := transaction(ss.store.DB, func(tx *sqlx.Tx) {
		_, txerr := tx.NamedExec(tables.Statements["upsertSurveyStructure"], survey)
		if txerr != nil {
			panic(txerr)
		}
		tx.MustExec(tables.Statements["updateAssignment"], survey.SAID)
	})
	return err
}

func (ss *SurveyStore) GetReport(surveyEventId int) ([]models.SurveyResult, error) {
	s := []models.SurveyResult{}
	err := ss.store.DB.Select(&s, tables.Statements["surveyReport"], surveyEventId)
	return s, err
}
