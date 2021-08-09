module github.com/HydrologicEngineeringCenter/nsi_survey_server

go 1.15

replace github.com/USACE/microauth v0.0.0 => /Users/rdcrlrsg/Projects/programming/go/src/github.com/usace/microauth

require (
	github.com/USACE/consequences-api v0.0.0-20201008020142-581a6decec1b
	github.com/USACE/microauth v0.0.0
	github.com/apex/gateway v1.1.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jmoiron/sqlx v1.3.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo/v4 v4.5.0
	github.com/lib/pq v1.10.2
	github.com/usace/dataquery v0.0.0-20210803012532-b16d2e028cb8 // indirect
)
