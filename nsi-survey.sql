create table survey_element (
    id serial primary key not null,
    fd_id int not null,
    survey_event_id int not null,
    is_control boolean
);

create table surveyor(
    user_id varchar(20) not null primary key,
    user_name varchar(200) not null
);

create table survey_assignment (
    id serial primary key not null,
    se_id int not null,
    completed boolean DEFAULT false,
    assigned_to varchar(20),
    CONSTRAINT fk_survey_element
        FOREIGN KEY(se_id) 
            REFERENCES survey_element(id),
    CONSTRAINT fk_user
        FOREIGN KEY(assigned_to) 
            REFERENCES surveyor(user_id)
);

create table survey_result(
    id serial not null primary key,
    sa_id int not null,
    fd_id int not null,
    X double precision not null,
    Y double precision not null,
    invalid_structure boolean not null,
    no_street_view boolean not null,
    cbfips varchar(15),
    occtype varchar(9),
    st_damcat varchar(3),
    found_ht double precision,
    num_story integer,
    sqft double precision,
    found_type varchar(4),
    rsmeans_type varchar(50),
    quality varchar(50),
    const_type varchar(50),
    garage varchar(50),
    roof_style varchar(50),

    CONSTRAINT fk_survey_assignment
        FOREIGN KEY(sa_id) 
            REFERENCES survey_assignment(id)
    
);

CREATE UNIQUE INDEX CONCURRENTLY idx_sr_said ON survey_result (sa_id);
ALTER TABLE survey_result ADD CONSTRAINT unique_sa_id UNIQUE USING INDEX idx_sr_said;

--drop table survey_result;
--drop table survey_assignment;
--drop table surveyor;
--drop table suvey_element;

/*
insert into surveyor values ('rr','Randy Goss');
insert into surveyor values ('ww','Will Lehman');
insert into surveyor values ('nn','Nick Lutz');
insert into surveyor values ('jj','Jack Goss');

insert into survey_element (fd_id,is_control) values (9,false);
insert into survey_element (fd_id,is_control) values (8,false);
insert into survey_element (fd_id,is_control) values (7,false);
insert into survey_element (fd_id,is_control) values (6,false);
insert into survey_element (fd_id,is_control) values (5,true);
insert into survey_element (fd_id,is_control) values (4,true);
insert into survey_element (fd_id,is_control) values (3,false);
insert into survey_element (fd_id,is_control) values (2,true);
insert into survey_element (fd_id,is_control) values (1,false);

insert into survey_assignment (se_id,assigned_to,completed) values (1,'nn',true);
insert into survey_assignment (se_id,assigned_to,completed) values (2,'rr',true);
insert into survey_assignment (se_id,assigned_to,completed) values (3,'rr',true);
insert into survey_assignment (se_id,assigned_to,completed) values (4,'ww',false);
insert into survey_assignment (se_id,assigned_to,completed) values (5,'rr',true);
insert into survey_assignment (se_id,assigned_to,completed) values (6,'rr',false);
insert into survey_assignment (se_id,assigned_to,completed) values (5,'nn',true);

select distinct
  t1.id as sa_id, 
  t1.se_id,
  t1.completed,
  (select (max(se_id)+1) from survey_assignment) as next_survey,
  (select min(t1.id) from survey_element t1 
      left outer join (select * from survey_assignment where assigned_to='nn') t2 on t1.id=t2.se_id
      where assigned_to is null and is_control='true') as next_control
from survey_element t2
left outer join survey_assignment t1 on t1.se_id=t2.id
where t1.id=(select max(id) from survey_assignment where assigned_to='nn') or t1.id is null
order by t1.id;

*/
