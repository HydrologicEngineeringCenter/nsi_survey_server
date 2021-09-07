package stores

import (
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
	dq "github.com/usace/goquery"
)

var surveyTable = dq.TableDataSet{
	Statements: map[string]string{
		"selectById": `select * from survey where id=$1`,
		"insert":     `insert into survey (title,description,active) values ($1,$2,$3) returning id`,
		"update":     `update survey set title=$1,description=$2,active=$3 where id=$4`,
		"nsi-survey": `select $2::uuid as sa_id, false as invalid_structure, false as no_street_view,fd_id,x,y,cbfips,occtype,st_damcat,found_ht,0.0 as num_story, 0.0 as sqft,found_type,
						'' as rsmeans_type, '' as quality, '' as const_type, '' as garage, '' as roof_style
						from nsi.nsi where fd_id=(select fd_id from survey_element where id=$1)`,
		"survey": `select sa_id, fd_id,x,y,invalid_structure,no_street_view,cbfips,occtype,st_damcat,found_ht,num_story,sqft,
					found_type,rsmeans_type,quality,const_type,garage,roof_style
					from survey_result where sa_id=$1`,
		"user-surveys": `select distinct s.id,s.title,s.description,s.active
							from survey s
							left outer join survey_owner so on so.survey_id=s.id
							left outer join survey_member sm on sm.survey_id=s.id
							where so.user_id=$1 or sm.user_id=$1`,
		"insert-owner": `insert into survey_member (survey_id,user_id,is_owner) values ($1,$2,$3)`,
	},
	Fields: models.Survey{},
}

var usersTable = dq.TableDataSet{
	Statements: map[string]string{
		"insert": `insert into users values ($1,$2)`,
	},
}

var surveyMemberTable = dq.TableDataSet{
	Statements: map[string]string{
		"upsert": `insert into survey_member(survey_id,user_id,is_owner) values ($1,$2,$3)
		                   ON CONFLICT(survey_id,user_id) do 
						  update set is_owner=EXCLUDED.is_owner`,
		"select_owners": "select * from survey_member where survey_id=$1",
		"remove":        `delete from survey_member where id=$1`,
	},
	Fields: models.SurveyMember{},
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
		"assignmentInfo": `select sa_id,se_id,completed,survey_order,next_survey_order,next_survey_seid,next_control_order,next_control_seid from (

								select t1.id as sa_id,t1.se_id,t1.completed,t2.survey_order, null as next_survey_order, null as next_survey_seid, null as next_control_order, null as next_control_seid
								from survey_element t2
								left outer join survey_assignment t1 on t1.se_id=t2.id
								where assigned_to=$1 and t2.survey_id=$2 and completed='false'
							
								union 
							
								select sa_id,se_id,completed,next_assignment.survey_order,next_survey_order, se1.id as next_survey_seid, next_control_order, se2.id as next_control_seid from 
									(select null::uuid as sa_id, null::uuid as se_id, null::bool as completed,null::integer as survey_order, 
							
									(select case when (
											select max(t2.survey_order) 
											from survey_assignment t1 
											inner join survey_element t2 on t2.id=t1.se_id 
											where t2.survey_id=$2 and t2.is_control='false') is null 
										then
											(select min(survey_order) from survey_element where survey_id=$2 and is_control='false')
										else
										(select min(survey_order) from survey_element where survey_order>
											(select max(t2.survey_order)
											from survey_assignment t1 
											inner join survey_element t2 on t2.id=t1.se_id 
											where t2.survey_id=$2 and t2.is_control='false') 
											and is_control='false')
										end) as next_survey_order,
									
									(select min(t1.survey_order) from survey_element t1
											left outer join (select * from survey_assignment where assigned_to=$1) t2 on t1.id=t2.se_id
											where assigned_to is null and is_control='true' and survey_id=$2) as next_control_order
								) next_assignment
								inner join survey_element se1 on se1.survey_order=next_assignment.next_survey_order
								left outer join survey_element se2 on se2.survey_order=next_assignment.next_control_order
								where se1.survey_id=$2 and (se2.survey_id=$2  or se2.survey_id is null)
							) assignment_query
	
							order by survey_order limit 1`,
	},
	Fields: models.SurveyAssignment{},
}

var resultTable = dq.TableDataSet{
	Statements: map[string]string{

		"nsi_survey": `select $2::uuid as sa_id, false as invalid_structure, false as no_street_view,fd_id,x,y,cbfips,occtype,st_damcat,found_ht,0 as num_story, 0.0 as sqft,found_type,
						'' as rsmeans_type, '' as quality, '' as const_type, '' as garage, '' as roof_style
						from nsi.nsi where fd_id=(select fd_id from survey_element where id=$1)`,

		"upsertSurveyStructure": `insert into survey_result
									(sa_id,fd_id,x,y,invalid_structure,no_street_view,cbfips,occtype,st_damcat,found_ht,num_story,sqft,found_type,rsmeans_type,quality,const_type,garage,roof_style)
									values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
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
