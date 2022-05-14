package gomologin

import (
	"database/sql"
	"fmt"
)

type DataBaseInterface interface {
	AuthenticateUser(username string, password string) (bool, interface{})
	RetriveRoles(username string) (bool, []string)
}

type SqlDataBase struct {
	*sql.DB
	AuthenticationSqlQuery string
	RolesSqlQuery          string
}

func (db *SqlDataBase) AuthenticateUser(username string, password string) (bool, interface{}) {
	var result interface{}
	err := db.QueryRow(db.AuthenticationSqlQuery, username, password).Scan(&result)

	if err != nil {
		fmt.Printf(err.Error())
	}

	return (err == nil), result

}

func (db *SqlDataBase) RetriveRoles(username string) (bool, []string) {

	if db.RolesSqlQuery == "" {
		return true, []string{}
	}

	rows, err := db.Query(db.RolesSqlQuery, username)

	if err != nil {
		fmt.Printf(err.Error())
		return true, []string{}
	}

	var roles []string
	for rows.Next() {
		var role string
		err := rows.Scan(&role)
		if err != nil {
			fmt.Println(err)
			return true, []string{}
		}

		roles = append(roles, role)
	}

	return (err == nil), roles

}
