package gologin

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"
)

// cookie handling
var Logintemplates = template.Must(template.ParseFiles(
	config.LoginPage,
))

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
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

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// login handler

func doLogin(response http.ResponseWriter, request *http.Request, db DataBaseInterface) {
	name := request.FormValue("username")
	pass := request.FormValue("password")
	redirectTarget := config.LoginPath + "?wrong"
	if name != "" && pass != "" {

		if name == "user" && pass == "1234" {
			setSession(name, response)
			redirectTarget = "/"
		}
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func loginView(w http.ResponseWriter, r *http.Request) {

	err := Logintemplates.ExecuteTemplate(w, "login", nil)

	if err != nil {
		http.Error(w, fmt.Sprintf("login: couldn't parse template: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			loginView(w, r)
		} else if r.Method == "POST" {
			var db DataBaseInterface
			if config.GetDBType() == "sql" {
				db = SqlDataBase{}
			}
			doLogin(w, r, db)
		}
	})
}

// logout handler
func LogoutHandler() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		clearSession(response)
		http.Redirect(response, request, "/", 302)
	})
}

func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userName := GetUserName(r)
		if userName != "" {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, config.LoginPath+"?expired", 302)
		}

	})
}
