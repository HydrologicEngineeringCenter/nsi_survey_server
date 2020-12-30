package models

type Config struct {
	SkipJWT       bool
	LambdaContext bool
	DBUser        string
	DBPass        string
	DBName        string
	DBHost        string
	DBSSLMode     string
	DBPort        string
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

type NsiStructure struct {
	FDID int `db:"fd_id"`
}
