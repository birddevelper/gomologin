package gologin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// func GetUserName(request *http.Request) (userName string) {
// 	if cookie, err := request.Cookie("session"); err == nil {
// 		cookieValue := make(map[string]string)
// 		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
// 			userName = cookieValue["name"]
// 		}
// 	}
// 	return userName
// }

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

	log.Printf("Set :" + sessionId)
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

func GetCurrentUsername(request *http.Request) string {
	username, ok := GetSession(username_session_key, request)
	if ok {
		return username.(string)
	} else {
		log.Println("getCurrentUsername error")
		return ""
	}
}

func GetDataReturnedByAuthQuery(request *http.Request) interface{} {
	data, ok := GetSession(auth_query_result_session_key, request)
	if ok {
		return data
	} else {
		log.Println("GetDataReturnedByAuthQuery error")
		return nil
	}
}

func doLogin(response http.ResponseWriter, request *http.Request, db DataBaseInterface) {
	username := request.FormValue("username")
	password := request.FormValue("password")
	redirectPath := request.URL.Query().Get("redirect")
	redirectTarget := config.LoginPath + "?wrong=yes&redirect=" + redirectPath
	if username != "" && password != "" {

		ok, data := db.AuthenticateUser(username, password)
		if ok {
			fmt.Printf("Underlying Value: %v\n", data)

			sessionId := generateSessionId(username)
			setSessionId(sessionId, response)
			setSessionBySessionId(sessionId, auth_query_result_session_key, data, request)
			setSessionBySessionId(sessionId, username_session_key, username, request)
			if redirectPath != "" {
				redirectTarget = redirectPath
			} else {
				redirectTarget = "/"
			}

		}

	}
	http.Redirect(response, request, redirectTarget, 302)
}

func loginView(response http.ResponseWriter, request *http.Request) {
	// cookie handling
	var Logintemplates = template.Must(template.ParseFiles(
		config.LoginPage,
	))
	err := Logintemplates.ExecuteTemplate(response, "login", nil)

	if err != nil {
		http.Error(response, fmt.Sprintf("login: couldn't parse template: %v", err), http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}

func LoginHandler() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		if request.Method == "GET" {
			loginView(response, request)
		} else if request.Method == "POST" {
			var db DataBaseInterface
			if config.GetDBType() == "sql" {
				db = &config.SqlDataBaseModel
			}
			doLogin(response, request, db)
		}
	})
}

// logout handler
func LogoutHandler() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		clearSession(response, request)
		http.Redirect(response, request, config.LoginPath+"?logout=yes", 302)
	})
}

func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		id := GetSessionId(request)
		if id != "" {
			next.ServeHTTP(response, request)
		} else {

			http.Redirect(response, request, config.LoginPath+"?redirect="+request.URL.Path, 302)
		}

	})
}
