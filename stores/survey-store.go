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

func (ss *SurveyStore) AddUser(user models.User) error {
	return ss.DS.Exec(goquery.NoTx, usersTable.Statements["insert"], user.UserID, user.Username)
}

func (ss *SurveyStore) GetSurveysforUser(userId string) (*[]models.Survey, error) {
	surveys := []models.Survey{}
	err := ss.DS.Select().
		DataSet(&surveyTable).
		StatementKey("user-surveys").
		Params(userId).
		Dest(&surveys).
		Fetch()
	return &surveys, err
}

func (ss *SurveyStore) GetSurveysforAdmin() (*[]models.Survey, error) {
	surveys := []models.Survey{}
	err := ss.DS.Select().
		DataSet(&surveyTable).
		StatementKey("admin-surveys").
		Params().
		Dest(&surveys).
		Fetch()
	return &surveys, err
}

func (ss *SurveyStore) GetSurveyMembers(surveyId uuid.UUID) (*[]models.SurveyMemberAlt, error) {
	members := []models.SurveyMemberAlt{}
	err := ss.DS.Select().
		DataSet(&surveyTable).
		StatementKey("members").
		Params(surveyId).
		Dest(&members).
		Fetch()
	return &members, err
}

func (ss *SurveyStore) GetSurveyElements(surveyId uuid.UUID) (*[]models.SurveyElementAlt, error) {
	elements := []models.SurveyElementAlt{}
	err := ss.DS.Select().
		DataSet(&surveyElementTable).
		StatementKey("select_elements").
		Params(surveyId).
		Dest(&elements).
		Fetch()
	return &elements, err
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
		_, err = ptx.Exec(context.Background(), surveyTable.Statements["insert-owner"], surveyId, userId, true)
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

func (ss *SurveyStore) UpsertSurveyMember(member models.SurveyMember) error {
	err := ss.DS.Exec(goquery.NoTx, surveyMemberTable.Statements["upsert"], member.SurveyID, member.UserID, member.IsOwner)
	return err
}

func (ss *SurveyStore) RemoveSurveyMember(memberId uuid.UUID) error {
	err := ss.DS.Exec(goquery.NoTx, surveyMemberTable.Statements["remove"], memberId)
	return err
}

func (ss *SurveyStore) RemoveMemberFromSurvey(memberId string, surveyId uuid.UUID) error {
	err := ss.DS.Exec(goquery.NoTx, surveyMemberTable.Statements["removeFromSurvey"], memberId, surveyId)
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
	err := ss.DS.Select(resultTable.Statements["surveyReport"]).
		Params(surveyId).
		Dest(&s).
		Fetch()
	return s, err
}

func (ss *SurveyStore) GetAssignmentInfo(userId string, surveyId uuid.UUID) (models.AssignmentInfo, error) {
	ai := models.AssignmentInfo{}
	err := ss.DS.Select(surveyAssignmentTable.Statements["assignmentInfo"]).
		Params(userId, surveyId).
		Dest(&ai).
		Fetch()

	if err != nil && err.Error() == NoResults {
		err = nil
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

func (ss *SurveyStore) SaveSurvey(survey *models.SurveyStructure) error {
	err := goquery.Transaction(ss.DS, func(tx goquery.Tx) {
		pgtx := tx.PgxTx()
		_, txerr := pgtx.Exec(context.Background(), resultTable.Statements["upsertSurveyStructure"],
			survey.SAID, survey.FDID, survey.X, survey.Y, survey.InvalidStructure, survey.NoStreetView,
			survey.CBfips, survey.OccupancyType, survey.Damcat, survey.FoundHt, survey.Stories, survey.SqFt,
			survey.FoundType, survey.RsmeansType, survey.Quality, survey.ConstType, survey.Garage, survey.RoofStyle)
		if txerr != nil {
			panic(txerr)
		}
		_, txerr = pgtx.Exec(context.Background(), surveyAssignmentTable.Statements["updateAssignment"], survey.SAID)
		if txerr != nil {
			panic(txerr)
		}
	})
	return err
}

func (ss *SurveyStore) IsOwner(surveyId uuid.UUID, userId string) bool {
	var owner int
	err := ss.DS.Select("select count(*) as owner from survey_member where survey_id=$1 and user_id=$2 and is_owner=true").
		Params(surveyId, userId).
		Dest(&owner).
		Fetch()
	if err != nil {
		log.Printf("Error in isOwner query:%s\n ", err)
		return false
	}
	return owner > 0
}

func (ss *SurveyStore) IsMember(surveyId uuid.UUID, userId string) bool {
	var member int
	err := ss.DS.Select("select count(*) as owner from survey_member where survey_id=$1 and user_id=$2").
		Params(surveyId, userId).
		Dest(&member).
		Fetch()
	if err != nil {
		log.Printf("Error in isMember query:%s\n ", err)
		return false
	}
	return member > 0
}
