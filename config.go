package gologin

import (
	"database/sql"
	"strings"
)

type Config struct {
	LoginPage         string
	SessionTimeout    int
	LoginPath         string
	SqlDataBase       *sql.DB
	AuthenticateQuery string
}

var config Config

func Configure() Config {
	config.LoginPage = "./templates/login.html"
	config.SessionTimeout = 120
	config.LoginPath = "/login"
	return config
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

func (config *Config) SqlDB(db *sql.DB) *Config {
	config.SqlDataBase = db
	return config
}

func (config *Config) SetAuthenticateQuery(query string) *Config {
	query = strings.Replace(query, "::username", "$1", 1)
	query = strings.Replace(query, "::password", "$2", 1)
	config.AuthenticateQuery = query
	return config
}

func (config *Config) GetDBType() string {
	if config.SqlDataBase != nil {
		return "sql"
	}
	return "noDB"
}
