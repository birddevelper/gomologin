package gologin

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"
)

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

func doLogin(response http.ResponseWriter, request *http.Request, db DataBaseInterface) {
	username := request.FormValue("username")
	password := request.FormValue("password")
	redirectPath := request.URL.Query().Get("redirect")
	redirectTarget := config.LoginPath + "?wrong=yes&redirect=" + redirectPath
	if username != "" && password != "" {

		ok, _ := db.AuthenticateUser(username, password)
		if ok {
			setSession(username, response)
			if redirectPath != "" {
				redirectTarget = redirectPath
			} else {
				redirectTarget = "/"
			}

		}

	}
	http.Redirect(response, request, redirectTarget, 302)
}

func loginView(w http.ResponseWriter, r *http.Request) {
	// cookie handling
	var Logintemplates = template.Must(template.ParseFiles(
		config.LoginPage,
	))
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
				db = &config.SqlDataBaseModel
			}
			doLogin(w, r, db)
		}
	})
}

// logout handler
func LogoutHandler() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		clearSession(response)
		http.Redirect(response, request, config.LoginPath+"?logout=yes", 302)
	})
}

func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userName := GetUserName(r)
		if userName != "" {
			next.ServeHTTP(w, r)
		} else {

			http.Redirect(w, r, config.LoginPath+"?redirect="+r.URL.Path, 302)
		}

	})
}
