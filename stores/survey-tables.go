package stores

import (
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
	dq "github.com/usace/goquery"
)

var surveyTable = dq.TableDataSet{
	Statements: map[string]string{
		"selectById":   `select * from survey where id=$1`,
		"insert":       `insert into survey (title,description,active) values ($1,$2,$3) returning id`,
		"insert-owner": `insert into survey_owner(survey_id,user_id) values ($1,$2)`,
		"update":       `update survey set title=$1,description=$2,active=$3 where id=$4`,
		"nsiSurvey": `select $2 as sa_id, false as invalid_structure, false as no_street_view,fd_id,x,y,cbfips,occtype,st_damcat,found_ht,0 as num_story, 0.0 as sqft,found_type,
						'' as rsmeans_type, '' as quality, '' as const_type, '' as garage, '' as roof_style
						from nsi.nsi where fd_id=(select fd_id from survey_element where id=$1)`,
		"survey": `select sa_id, fd_id,x,y,invalid_structure,no_street_view,cbfips,occtype,st_damcat,found_ht,num_story,sqft,
					found_type,rsmeans_type,quality,const_type,garage,roof_style
					from survey_result where sa_id=$1`,
	},
	Fields: models.Survey{},
}

var surveyOwnerTable = dq.TableDataSet{
	Statements: map[string]string{
		"insert":        `insert into survey_owner(survey_id,user_id) values ($1,$2)`,
		"select_owners": "select * from survey_owner where survey_id=$1",
		"remove":        `delete from survey_owner where id=$1`,
	},
	Fields: models.SurveyOwner{},
}

var surveyElementTable = dq.TableDataSet{
	Name:   "survey_element",
	Fields: models.SurveyElement{},
}

var surveyAssignmentTable = dq.TableDataSet{
	Name: "survey_assignment",
	Statements: map[string]string{
		"updateAssignment": `update survey_assignment set completed='true' where id=$1`,
		"assignSurvey":     `insert into survey_assignment (se_id,assigned_to) values ($1,$2) returning id`,
		"assignmentInfo": `select distinct
			t1.id as sa_id,
			t1.se_id,
			t1.completed,
			(select (max(t1.se_id)+1) from survey_assignment t1 inner join survey_element t2 on t2.id=t1.se_id where t2.survey_event_id=$1) as next_survey,
			(select min(t1.id) from survey_element t1
					left outer join (select * from survey_assignment where assigned_to=$2) t2 on t1.id=t2.se_id
					where assigned_to is null and is_control='true' and survey_event_id=$3) as next_control
			from survey_element t2
			left outer join survey_assignment t1 on t1.se_id=t2.id
			where t1.id=(select max(t1.id) from survey_assignment t1 inner join survey_element t2 on t1.se_id=t2.id where t2.survey_event_id=$4 and t1.assigned_to=$5)
				or t1.id is null
			order by t1.id`,
	},
	Fields: models.SurveyAssignment{},
}

var miscQueries = dq.TableDataSet{
	Statements: map[string]string{

		"nsi_survey": `select $2 as sa_id, false as invalid_structure, false as no_street_view,fd_id,x,y,cbfips,occtype,st_damcat,found_ht,0 as num_story, 0.0 as sqft,found_type,
						'' as rsmeans_type, '' as quality, '' as const_type, '' as garage, '' as roof_style
						from nsi.nsi where fd_id=(select fd_id from survey_element where id=$1)`,

		"upsertSurveyStructure": `insert into survey_result
									(sa_id,fd_id,x,y,invalid_structure,no_street_view,cbfips,occtype,st_damcat,found_ht,num_story,sqft,found_type,rsmeans_type,quality,const_type,garage,roof_style)
									values (:sa_id,:fd_id,:x,:y,:invalid_structure,:no_street_view,:cbfips,:occtype,:st_damcat,:found_ht,:num_story,:sqft,:found_type,:rsmeans_type,:quality,:const_type,:garage,:roof_style)
									ON CONFLICT (sa_id)
									DO UPDATE SET x=EXCLUDED.x,y=EXCLUDED.y,invalid_structure=EXCLUDED.invalid_structure,no_street_view=EXCLUDED.no_street_view, cbfips=EXCLUDED.cbfips,
													occtype=EXCLUDED.occtype,st_damcat=EXCLUDED.st_damcat,found_ht=EXCLUDED.found_ht,num_story=EXCLUDED.num_story,
												sqft=EXCLUDED.sqft,found_type=EXCLUDED.found_type,rsmeans_type=EXCLUDED.rsmeans_type,
												quality=EXCLUDED.quality,const_type=EXCLUDED.const_type,garage=EXCLUDED.garage,roof_style=EXCLUDED.roof_style`,

		"surveyReport": `select
				t1.id as sr_id,
				t3.user_id,
				t3.user_name,
				t2.completed,
				t4.is_control,
				t1.sa_id,
				t1.fd_id,
				t1.x,
				t1.y,
				t1.cbfips,
				t1.occtype,
				t1.st_damcat,
				t1.found_ht,
				t1.num_story,
				t1.sqft,
				t1.found_type,
				t1.rsmeans_type,
				t1.quality,
				t1.const_type,
				t1.garage,
				t1.roof_style,
				t1.invalid_structure,
				t1.no_street_view
				from survey_result t1
				inner join survey_assignment t2 on t2.id=t1.sa_id
				inner join surveyor t3 on t3.user_id=t2.assigned_to
				inner join survey_element t4 on t4.id=t2.se_id
				where t4.survey_event_id=$1`,
	},
}
