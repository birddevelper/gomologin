package gologin

import "fmt"

var session_store map[string]map[string]interface{}
var username_session_key string = "current_username_"
var auth_query_result_session_key = "auth_query_result_"

func init() {

	// initiate empty session store
	session_store = make(map[string]map[string]interface{})
	username_session_key = username_session_key + currentTimeStamp()
	auth_query_result_session_key = auth_query_result_session_key + currentTimeStamp()
	fmt.Println("Session store initialized!")
}
