package gologin

import (
	"database/sql"
	"strings"
)

type Config struct {
	LoginPage        string
	SessionTimeout   int
	LoginPath        string
	SqlDataBaseModel SqlDataBase
}

var config Config

func Configure() *Config {

	config.LoginPage = "./templates/login.html"
	config.SessionTimeout = 120
	config.LoginPath = "/login"
	return &config
}

func (config *Config) SetLoginPage(loginPage string) *Config {
	config.LoginPage = loginPage
	return config
}

func (config *Config) SetSessionTimeout(sessionTimeout int) *Config {
	config.SessionTimeout = sessionTimeout
	return config
}

func (config *Config) SetLoginPath(loginPath string) *Config {
	config.LoginPath = loginPath
	return config
}

func (config *Config) AuthenticateBySqlQuery(db *sql.DB, authenticateQuery string, rolesQuery string) *Config {

	authenticateQuery = strings.Replace(authenticateQuery, "::username", "?", 1)
	authenticateQuery = strings.Replace(authenticateQuery, "::password", "?", 1)

	rolesQuery = strings.Replace(rolesQuery, "::username", "?", 1)

	config.SqlDataBaseModel = SqlDataBase{db, authenticateQuery, rolesQuery}
	return config
}

func (config *Config) GetDBType() string {
	if config.SqlDataBaseModel != (SqlDataBase{}) {
		return "sql"
	}
	return "noDB"
}
