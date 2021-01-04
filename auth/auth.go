package auth

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/HydrologicEngineeringCenter/nsi_survey_server/models"
	"github.com/HydrologicEngineeringCenter/nsi_survey_server/stores"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	PUBLIC = iota
	ORGADMIN
)

type Auth struct {
	Store     *stores.SurveyStore
	VerifyKey *rsa.PublicKey
}

func (a *Auth) Authorize(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		claims, err := a.marshalJwt(tokenString)
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusUnauthorized, "bad token")
		}
		c.Set("NSIUSER", claims)
		return handler(c)
	}
}

func (a *Auth) LoadVerificationKey(filePath string) error {
	publicKeyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	pk, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return err
	}
	a.VerifyKey = pk
	return nil
}

func (a *Auth) marshalJwt(tokenString string) (models.JwtClaim, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.VerifyKey, nil
	})
	if err != nil {
		return models.JwtClaim{}, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jwtUser := models.JwtClaim{
			Sub:  claims["sub"].(string),
			Name: claims["name"].(string),
		}
		return jwtUser, nil
	} else {
		return models.JwtClaim{}, errors.New("Invalid Token")
	}

}

func contains(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
