package stores

import (
	"fmt"
	"log"
	"strconv"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type SurveyStore struct {
	db *sqlx.DB
}

func CreateSurveyStore(appConfig *models.Config) (*SurveyStore, error) {
	dburl := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		appConfig.Dbuser, appConfig.Dbpass, appConfig.Dbhost, appConfig.Dbport, appConfig.Dbname)
	con, err := sqlx.Connect("pgx", dburl)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected as %s to database %s:%s/%s", appConfig.Dbuser, appConfig.Dbhost, appConfig.Dbport, appConfig.Dbname)
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

var surveySql string = `select $2 as sa_id, fd_id,x,y,cbfips 
                        from nsi.nsi where fd_id=(select fd_id from survey_element where id=$1)`

func (ss *SurveyStore) GetStructure(seId int, saId int) (models.SurveyStructure, error) {
	s := models.SurveyStructure{}
	err := ss.db.Get(&s, surveySql, seId, strconv.Itoa(saId))
	if err != nil {
		log.Printf("Failed to retrieve structure: %s/n", err)
		return s, err
	}
	return s, nil
}

var assignSurveySql string = `insert into survey_assignment (se_id,assigned_to) values ($1,$2)`

func (ss *SurveyStore) AssignSurvey(userId string, seId int) (int, error) {
	res, err := ss.db.Exec(assignSurveySql, seId, userId)
	if err != nil {
		return -1, err
	}

	saId, sa_err := res.LastInsertId()
	if sa_err != nil {
		return -1, sa_err
	}

	return int(saId), nil
}

var insertSurveyStructure string = `insert into survey_structure (sa_id,fd_id,x,y,cbfips) values (:sa_id,:fd_id,:x,:y,:cbfips)`
var updateAssignment string = `update survey_assignment set completed='true' where sa_id=$1`

func (ss *SurveyStore) SaveSurvey(survey *models.SurveyStructure) error {
	err := transaction(ss.db, func(tx *sqlx.Tx) {
		_, txerr := tx.NamedExec(insertSurveyStructure, survey)
		if txerr != nil {
			log.Panicf("Unable to insert survey: %s", txerr)
		}
		tx.MustExec(insertSurveyStructure, survey.SAID)
	})
	return err
}
