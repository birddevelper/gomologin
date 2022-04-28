package gologin

import "fmt"

var session_store map[string]map[string]interface{}

func init() {

	// initiate empty session store
	session_store = make(map[string]map[string]interface{})

	fmt.Println("Session store initialized!")
}
