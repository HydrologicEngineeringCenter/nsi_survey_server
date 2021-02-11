package stores

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

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
	con.SetMaxOpenConns(4)

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
	if err != nil {
		return models.AssignmentInfo{}, err
	}
	if len(ai) == 0 {
		return models.AssignmentInfo{}, errors.New("Invalid Record")
	}
	if ai[0].NextSurvey == nil {
		ns := 1
		ai[0].NextSurvey = &ns
	}
	return ai[0], err
}

var nsiSurveySql string = `select $2 as sa_id, false as invalid_structure, fd_id,x,y,cbfips,occtype,st_damcat,found_ht,0 as num_story, 0.0 as sqft,found_type,
                        '' as rsmeans_type, '' as quality, '' as const_type, '' as garage, '' as roof_style 
						from nsi.nsi where fd_id=(select fd_id from survey_element where id=$1)`

var surveySql string = `select sa_id, fd_id,x,y,invalid_structure,cbfips,occtype,st_damcat,found_ht,num_story,sqft,
                        found_type,rsmeans_type,quality,const_type,garage,roof_style 
                        from survey_result where sa_id=$1`

func (ss *SurveyStore) GetStructure(seId int, saId int) (models.SurveyStructure, error) {
	s := models.SurveyStructure{}
	err := ss.db.Get(&s, surveySql, strconv.Itoa(saId))
	if err != nil {
		if err == sql.ErrNoRows {
			//no existing survey result, get survey data from nsi
			err := ss.db.Get(&s, nsiSurveySql, seId, strconv.Itoa(saId))
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

var assignSurveySql string = `insert into survey_assignment (se_id,assigned_to) values ($1,$2) returning id`

func (ss *SurveyStore) AssignSurvey(userId string, seId int) (int, error) {
	var saId int
	err := ss.db.QueryRow(assignSurveySql, seId, userId).Scan(&saId)
	if err != nil {
		return -1, err
	}
	return int(saId), nil
}

/*
var insertSurveyStructure string = `insert into survey_result
									 (sa_id,fd_id,x,y,cbfips,occtype,st_damcat,found_ht,num_story,sqft,found_type,
									  rsmeans_type,quality,const_type,garage,roof_style)
									 values (:sa_id,:fd_id,:x,:y,:cbfips,:occtype,:st_damcat,:found_ht,:num_story,:sqft,:found_type,
									  :rsmeans_type,:quality,:const_type,:garage,:roof_style)`
*/

var insertSurveyStructure string = `insert into survey_result 
									  (sa_id,fd_id,x,y,invalid_structure,cbfips,occtype,st_damcat,found_ht,num_story,sqft,found_type,rsmeans_type,quality,const_type,garage,roof_style) 
									  values (:sa_id,:fd_id,:x,:y,:invalid_structure,:cbfips,:occtype,:st_damcat,:found_ht,:num_story,:sqft,:found_type,:rsmeans_type,:quality,:const_type,:garage,:roof_style)
									  ON CONFLICT (sa_id)
									  DO UPDATE SET x=EXCLUDED.x,y=EXCLUDED.y,invalid_structure=EXCLUDED.invalid_structure,cbfips=EXCLUDED.cbfips,occtype=EXCLUDED.occtype,
													st_damcat=EXCLUDED.st_damcat,found_ht=EXCLUDED.found_ht,num_story=EXCLUDED.num_story,
													sqft=EXCLUDED.sqft,found_type=EXCLUDED.found_type,rsmeans_type=EXCLUDED.rsmeans_type,
													quality=EXCLUDED.quality,const_type=EXCLUDED.const_type,garage=EXCLUDED.garage,roof_style=EXCLUDED.roof_style`

var updateAssignment string = `update survey_assignment set completed='true' where id=$1`

func (ss *SurveyStore) SaveSurvey(survey *models.SurveyStructure) error {
	err := transaction(ss.db, func(tx *sqlx.Tx) {
		_, txerr := tx.NamedExec(insertSurveyStructure, survey)
		if txerr != nil {
			panic(txerr)
		}
		tx.MustExec(updateAssignment, survey.SAID)
	})
	return err
}
