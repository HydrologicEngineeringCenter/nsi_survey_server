package config

import dq "github.com/usace/goquery"

type Config struct {
	SkipJWT       bool
	LambdaContext bool
	Dbuser        string
	Dbpass        string
	Dbname        string
	Dbhost        string
	Dbstore       string
	Dbdriver      string
	DBSSLMode     string
	Dbport        string
	Ippk          string
	Port          string
	Aud           string
}

func (c *Config) Rdbmsconfig() dq.RdbmsConfig {
	return dq.RdbmsConfig{
		Dbuser:   c.Dbuser,
		Dbpass:   c.Dbpass,
		Dbhost:   c.Dbhost,
		Dbport:   c.Dbport,
		Dbname:   c.Dbname,
		DbDriver: c.Dbdriver,
		DbStore:  c.Dbstore,
	}
}
