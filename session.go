package gomologin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
)

var session_store map[string]map[string]interface{}
var username_session_key string = "current_username_"
var auth_query_result_session_key = "auth_query_result_"
var roles_session_key = "current_user_roles"

func init() {

	// initiate empty session store
	session_store = make(map[string]map[string]interface{})
	username_session_key = username_session_key + currentTimeStamp()
	auth_query_result_session_key = auth_query_result_session_key + currentTimeStamp()
	roles_session_key = roles_session_key + currentTimeStamp()
	fmt.Println("Session store initialized!")
}

func generateSessionId(username string) string {
	//now := time.Now()
	//unix := now.Unix()
	timestamp := currentTimeStamp() //strconv.FormatInt(unix, 10)
	hash := md5.Sum([]byte(username + timestamp))
	return hex.EncodeToString(hash[:])
}

func setSessionId(id string, response http.ResponseWriter) {

	session_store[id] = make(map[string]interface{})

	value := map[string]string{
		"id": id,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:   "session",
			Value:  encoded,
			Path:   "/",
			MaxAge: config.SessionTimeout,
		}
		http.SetCookie(response, cookie)
	}
}

func GetSessionId(request *http.Request) (id string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			id = cookieValue["id"]
		}
	}
	return id
}

func SetSession(key string, value interface{}, request *http.Request) {
	sessionId := GetSessionId(request)
	session_store[sessionId][key] = value

}

func setSessionBySessionId(sessionId string, key string, value interface{}, request *http.Request) {
	session_store[sessionId][key] = value

}

func GetSession(key string, request *http.Request) (interface{}, bool) {
	sessionId := GetSessionId(request)
	value, ok := session_store[sessionId][key]
	if ok {
		return value, true
	}

	return nil, false
}

func RemoveSession(key string, request *http.Request) bool {

	sessionId := GetSessionId(request)

	if _, ok := session_store[sessionId][key]; ok {
		delete(session_store[sessionId], key)
		return true
	}

	return false
}

func clearSession(response http.ResponseWriter, request *http.Request) {

	sessionId := GetSessionId(request)
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	if _, ok := session_store[sessionId]; ok {
		delete(session_store, sessionId)
	}

	http.SetCookie(response, cookie)
}
