package auth

import (
	"github.com/USACE/microauth"
	"github.com/labstack/echo/v4"
)

const (
	PUBLIC = iota
	ORGADMIN
)

func Appauth(c echo.Context, claims microauth.JwtClaim) bool {
	c.Set("NSIUSER", claims)
	return true
}
