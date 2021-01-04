package models

type JwtClaim struct {
	Sub  string
	Name string
}

type Config struct {
	SkipJWT       bool
	LambdaContext bool
	Dbuser        string
	Dbpass        string
	Dbname        string
	Dbhost        string
	DBSSLMode     string
	Dbport        string
	Ippk          string
}

type AssignmentInfo struct {
	SA_ID       *int  `db:"sa_id"`
	SE_ID       *int  `db:"se_id"`
	Completed   *bool `db:"completed"`
	NextSurvey  int   `db:"next_survey"`
	NextControl int   `db:"next_control"`
}

type SurveyAssignment struct {
	ID               string
	SurveyElement_ID string
	Completed        bool
	Assigned         string
}

type SurveyElement struct {
	ID         string `json:"se_id"`
	FD_ID      string `json:"fd_id"`
	Is_control bool   `json:"is_control"`
}

type SurveyResult struct {
	ID    string `json:"sr_id"`
	FD_ID string `json:"fd_id"`
}

type SurveyStructure struct {
	SAID          int     `db:"sa_id" json:"saId"`
	FDID          int     `db:"fd_id" json:"fdId"`
	X             float64 `db:"x" json:"x"`
	Y             float64 `db:"y" json:"y"`
	CBfips        string  `db:"cbfips" json:"cbfips"`
	OccupancyType string  `db:"occtype" json:"occupancyType"`
	Damcat        string  `db:"st_damcat" json:"damcat"`
	FoundHt       float64 `db:"found_ht" json:"found_ht"`
	Stories       int     `db:"num_story" json:"stories"`
	SqFt          float64 `db:"sqft" json:"sq_ft"`
	FoundType     string  `db:"found_type" json:"found_type"`
	RsmeansType   string  `db:"rsmeans_type" json:"rsmeans_type"`
	Quality       string  `db:"quality" json:"quality"`
	ConstType     string  `db:"const_type" json:"const_type"`
	Garage        string  `db:"garage" json:"garage"`
	RoofStyle     string  `db:"roof_style" json:"roof_style"`
}
