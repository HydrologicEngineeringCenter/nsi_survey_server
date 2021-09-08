module github.com/HydrologicEngineeringCenter/nsi_survey_server

go 1.15

replace github.com/USACE/microauth v0.0.0 => /Users/rdcrlrsg/Projects/programming/go/src/github.com/usace/microauth

replace github.com/usace/goquery v0.0.0-20210825130028-45f9bdbb50fa => /Users/rdcrlrsg/Projects/programming/go/src/github.com/usace/goquery

require (
	github.com/USACE/consequences-api v0.0.0-20201008020142-581a6decec1b
	github.com/USACE/microauth v0.0.0
	github.com/apex/gateway v1.1.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/jmoiron/sqlx v1.3.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo/v4 v4.5.0
	github.com/lib/pq v1.10.2
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/usace/goquery v0.0.0-20210825130028-45f9bdbb50fa // indirect
	golang.org/x/tools v0.1.5 // indirect
// indirect
)
