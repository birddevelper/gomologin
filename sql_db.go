package gologin

type DataBaseInterface interface {
	AuthenticateUser(username string, password string) (bool, interface{})
	RetriveRoles(username string) (bool, []string)
}

type SqlDataBase struct {
}

func (db SqlDataBase) AuthenticateUser(username string, password string) (bool, interface{}) {
	var result interface{}
	err := config.SqlDataBase.QueryRow(config.AuthenticateQuery, username, password).Scan(&result)

	return (err == nil), result

}

func (db SqlDataBase) RetriveRoles(username string) (bool, []string) {

	return true, []string{"ADMIN", "USER"}

}
