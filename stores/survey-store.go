package stores

import (
	"fmt"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
	"github.com/jmoiron/sqlx"
)

type SurveyStore struct {
	db *sqlx.DB
}

func CreateSurveyStore(appConfig *models.Config) (*SurveyStore, error) {
	dburl := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		appConfig.DBUser, appConfig.DBPass, appConfig.DBHost, appConfig.DBPort, appConfig.DBName)
	con, err := sqlx.Connect("pgx", dburl)
	if err != nil {
		return nil, err
	}
	con.SetMaxOpenConns(10)

	ss := SurveyStore{
		db: con,
	}
	return &ss, nil
}

var assignmentInfoSql string = `select distinct
								t1.id as sa_id, 
								t1.se_id,
								t1.completed,
								(select (max(se_id)+1) from survey_assignment) as next_survey,
								(select min(t1.id) from survey_element t1 
									left outer join (select * from survey_assignment where assigned_to=$1) t2 on t1.id=t2.se_id
									where assigned_to is null and is_control='true') as next_control
								from survey_element t2
								left outer join survey_assignment t1 on t1.se_id=t2.id
								where t1.id=(select max(id) from survey_assignment where assigned_to=$2) or t1.id is null
								order by t1.id`

func (ss *SurveyStore) GetAssignmentInfo(userId string) (models.AssignmentInfo, error) {
	ai := []models.AssignmentInfo{}
	err := ss.db.Select(&ai, assignmentInfoSql, userId, userId)
	return ai[0], err
}

var surveySql string = `select * from nsi.nsi where fd_id=(select fd_id from survey_elements where se_id=$1)`

func (ss *SurveyStore) GetStructure(seId int) (models.NsiStructure, error) {
	s := models.NsiStructure{}
	err := ss.db.Get(&s, surveySql, seId)
	if err != nil {
		return s, err
	}
	return s, nil
}

var surveyAssignInsertSql string = `insert into survey_assignment (se_id,assigned_to) values ($1,$2)`

func (ss *SurveyStore) AssignSurvey(userId string, seId int) error {
	_, err := ss.db.Exec(surveyAssignInsertSql, seId, userId)
	return err
}
