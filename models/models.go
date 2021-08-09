package models

import (
	"fmt"
	"strconv"

	"github.com/usace/dataquery"
)

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
	SurveyEvent   int
}

func (c *Config) Rdbmsconfig() dataquery.RdbmsConfig {
	return dataquery.RdbmsConfig{
		Dbuser: c.Dbuser,
		Dbpass: c.Dbpass,
		Dbhost: c.Dbhost,
		Dbport: c.Dbport,
		Dbname: c.Dbname,
	}
}

type AssignmentInfo struct {
	SA_ID       *int  `db:"sa_id"`
	SE_ID       *int  `db:"se_id"`
	Completed   *bool `db:"completed"`
	NextSurvey  *int  `db:"next_survey"`
	NextControl *int  `db:"next_control"`
}

type SurveyAssignment struct {
	ID               string
	SurveyElement_ID string
	Completed        bool
	Assigned         string
}

type SurveyElement struct {
	ID            string `json:"se_id"`
	FD_ID         string `json:"fd_id"`
	Is_control    bool   `json:"is_control"`
	SurveyEventID int    `json:"surveyEventId"`
}

/*
type SurveyResult struct {
	ID    string `json:"sr_id"`
	FD_ID string `json:"fd_id"`
}
*/

type SurveyStructure struct {
	SAID             int     `db:"sa_id" json:"saId"`
	FDID             int     `db:"fd_id" json:"fdId"`
	X                float64 `db:"x" json:"x"`
	Y                float64 `db:"y" json:"y"`
	InvalidStructure bool    `db:"invalid_structure" json:"invalidStructure"`
	NoStreetView     bool    `db:"no_street_view" json:"noStreetView"`
	CBfips           string  `db:"cbfips" json:"cbfips"`
	OccupancyType    string  `db:"occtype" json:"occupancyType"`
	Damcat           string  `db:"st_damcat" json:"damcat"`
	FoundHt          float64 `db:"found_ht" json:"found_ht"`
	Stories          float64 `db:"num_story" json:"stories"`
	SqFt             float64 `db:"sqft" json:"sq_ft"`
	FoundType        string  `db:"found_type" json:"found_type"`
	RsmeansType      string  `db:"rsmeans_type" json:"rsmeans_type"`
	Quality          string  `db:"quality" json:"quality"`
	ConstType        string  `db:"const_type" json:"const_type"`
	Garage           string  `db:"garage" json:"garage"`
	RoofStyle        string  `db:"roof_style" json:"roof_style"`
}

type SurveyResult struct {
	SRID      int    `db:"sr_id" json:"srId"`
	UserID    string `db:"user_id" json:"userId"`
	UserName  string `db:"user_name" json:"userName"`
	Completed bool   `db:"completed" json:"completed"`
	IsControl bool   `db:"is_control" json:"isControl"`

	SurveyStructure
}

func (sr SurveyResult) String() []string {
	return ([]string{
		strconv.Itoa(sr.SRID),
		fmt.Sprintf(`"%s"`, sr.UserID),
		fmt.Sprintf(`"%s"`, sr.UserName),
		strconv.FormatBool(sr.Completed),
		strconv.FormatBool(sr.IsControl),
		strconv.Itoa(sr.SAID),
		strconv.Itoa(sr.FDID),
		strconv.FormatFloat(sr.X, 'f', 8, 64),
		strconv.FormatFloat(sr.Y, 'f', 8, 64),
		strconv.FormatBool(sr.InvalidStructure),
		strconv.FormatBool(sr.NoStreetView),
		fmt.Sprintf(`"%s"`, sr.CBfips),
		fmt.Sprintf(`"%s"`, sr.OccupancyType),
		fmt.Sprintf(`"%s"`, sr.Damcat),
		strconv.FormatFloat(sr.FoundHt, 'f', 4, 64),
		strconv.FormatFloat(sr.Stories, 'f', 4, 64),
		strconv.FormatFloat(sr.SqFt, 'f', 4, 64),
		fmt.Sprintf(`"%s"`, sr.FoundType),
		fmt.Sprintf(`"%s"`, sr.RsmeansType),
		fmt.Sprintf(`"%s"`, sr.Quality),
		fmt.Sprintf(`"%s"`, sr.ConstType),
		fmt.Sprintf(`"%s"`, sr.Garage),
		fmt.Sprintf(`"%s"`, sr.RoofStyle),
	})
}
