package models

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type JwtClaim struct {
	Sub  string
	Name string
}

type Survey struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Active      bool      `db:"active" json:"active"`
}

type User struct {
	UserID   string `db:"user_id" json:"userId"`
	Username string `db:"user_name" json:"userName"`
}
type SurveyMember struct {
	ID       uuid.UUID `db:"id" json:"id"`
	SurveyID uuid.UUID `db:"survey_id" json:"surveyId"`
	UserID   string    `db:"user_id" json:"userId"`
	IsOwner  bool      `db:"is_owner" json:"isOwner"`
}

type AssignmentInfo struct {
	SAID             *uuid.UUID `db:"sa_id"`
	SEID             *uuid.UUID `db:"se_id"`
	Completed        *bool      `db:"completed"`
	SurveyOrder      *int       `db:"survey_order"`
	NextSurveyOrder  *int       `db:"next_survey_order"`
	NextSurveySEID   *uuid.UUID `db:"next_survey_seid"`
	NextControlOrder *int       `db:"next_control_order"`
	NextControlSEID  uuid.UUID  `db:"next_control_seid"`
}

type SurveyElement struct {
	ID          uuid.UUID `json:"seId" db:"id" dbid:"AUTOINCREMENT"`
	SurveyID    uuid.UUID `json:"surveyId" db:"survey_id"`
	SurveyOrder int       `json:"surveyOrder" db:"survey_order"`
	FD_ID       int       `json:"fdId" db:"fd_id"`
	Is_control  bool      `json:"isControl" db:"is_control"`
}

type SurveyAssignment struct {
	ID               uuid.UUID `json:"saId" db:"id" dbid:"AUTOINCREMENT"`
	SurveyElement_ID uuid.UUID `json:"seId" db:"se_id"`
	Completed        bool      `json:"completed" db:"completed"`
	Assigned         string    `json:"assignedTo" db:"assigned_to"`
}

type SurveyStructure struct {
	SAID             uuid.UUID `db:"sa_id" json:"saId"`
	FDID             int       `db:"fd_id" json:"fdId"`
	X                float64   `db:"x" json:"x"`
	Y                float64   `db:"y" json:"y"`
	InvalidStructure bool      `db:"invalid_structure" json:"invalidStructure"`
	NoStreetView     bool      `db:"no_street_view" json:"noStreetView"`
	CBfips           string    `db:"cbfips" json:"cbfips"`
	OccupancyType    string    `db:"occtype" json:"occupancyType"`
	Damcat           string    `db:"st_damcat" json:"damcat"`
	FoundHt          float64   `db:"found_ht" json:"found_ht"`
	Stories          int       `db:"num_story" json:"stories"`
	SqFt             float64   `db:"sqft" json:"sq_ft"`
	FoundType        string    `db:"found_type" json:"found_type"`
	RsmeansType      string    `db:"rsmeans_type" json:"rsmeans_type"`
	Quality          string    `db:"quality" json:"quality"`
	ConstType        string    `db:"const_type" json:"const_type"`
	Garage           string    `db:"garage" json:"garage"`
	RoofStyle        string    `db:"roof_style" json:"roof_style"`
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
		sr.SAID.String(),
		strconv.Itoa(sr.FDID),
		strconv.FormatFloat(sr.X, 'f', 8, 64),
		strconv.FormatFloat(sr.Y, 'f', 8, 64),
		strconv.FormatBool(sr.InvalidStructure),
		strconv.FormatBool(sr.NoStreetView),
		fmt.Sprintf(`"%s"`, sr.CBfips),
		fmt.Sprintf(`"%s"`, sr.OccupancyType),
		fmt.Sprintf(`"%s"`, sr.Damcat),
		strconv.FormatFloat(sr.FoundHt, 'f', 4, 64),
		strconv.Itoa(sr.Stories),
		strconv.FormatFloat(sr.SqFt, 'f', 4, 64),
		fmt.Sprintf(`"%s"`, sr.FoundType),
		fmt.Sprintf(`"%s"`, sr.RsmeansType),
		fmt.Sprintf(`"%s"`, sr.Quality),
		fmt.Sprintf(`"%s"`, sr.ConstType),
		fmt.Sprintf(`"%s"`, sr.Garage),
		fmt.Sprintf(`"%s"`, sr.RoofStyle),
	})
}
