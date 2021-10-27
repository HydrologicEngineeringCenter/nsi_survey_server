package auth

import (
	"log"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/stores"
	. "github.com/USACE/microauth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	PUBLIC = iota
	ADMIN
	SURVEY_OWNER
	SURVEY_MEMBER
)

func Appauth(c echo.Context, authstore interface{}, roles []int, claims JwtClaim) bool {
	c.Set("NSIUSER", claims)
	store := authstore.(*stores.SurveyStore)
	store.AddUser(models.User{
		UserID:   claims.Sub,
		Username: claims.UserName,
	})

	surveyId, err := uuid.Parse(c.Param("surveyid"))

	if err != nil { // even with no surveyid, there are handlers that can be invoked
		log.Printf("Invalid survey_id: %s\n", err)
		// return false

	} else { // surveyId exists in URI

		c.Set("NSISURVEY", surveyId)

		if Contains(roles, SURVEY_OWNER) {
			return store.IsOwner(surveyId, claims.Sub)
		}
		if Contains(roles, SURVEY_MEMBER) {
			return store.IsMember(surveyId, claims.Sub)
		}
	}

	if Contains(roles, PUBLIC) {
		return true
	}
	if Contains(roles, ADMIN) && Contains_string(claims.Roles, "ADMIN") {
		return true
	}

	return false
}

// func containsParam(paramNames []string, param string) bool {
// 	for _, pn := range paramNames {
// 		if pn == param {
// 			return true
// 		}
// 	}
// 	return false
// }
