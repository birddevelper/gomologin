package gologin

import "database/sql"

type DataBaseInterface interface {
	AuthenticateUser(username string, password string) (bool, interface{})
	RetriveRoles(username string) (bool, []string)
}

type SqlDataBase struct {
	*sql.DB
	AuthenticationSqlQuery string
}

func (db *SqlDataBase) AuthenticateUser(username string, password string) (bool, interface{}) {
	var result interface{}
	err := db.QueryRow(db.AuthenticationSqlQuery, username, password).Scan(&result)

	return (err == nil), result

}

func (db *SqlDataBase) RetriveRoles(username string) (bool, []string) {

	return true, []string{"ADMIN", "USER"}

}
