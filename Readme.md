# What is Gologin

<p align="center">
<img src="https://m-shaeri.ir/blog/wp-content/uploads/2022/04/gologin.png"  height="200" >
</p>

**Gologin** is an easy to setup login manager for Go web applications. It helps you protect your application resources from unattended/unauthenticated access. Currently it works with SQL databases authentication.

## How to setup

Get the package with this command :

```bash
go get github.com/birddevelper/gologin

```

## How to use

You can easily setup your customized login process with **configure()** function. You should specify following paramters to make the Gologin ready to start:

- **Login page** : path to html template. Default path is ***./template/login.html***, note that the template must be defined as ****"login"**** with ***{{define "login"}}*** at the begining line

- **Login path** : login http path. Default path is ***/login***

- **Session timeout** : Number of seconds before the session expires. Default value is 120 seconds.

- **SQL connection and query to authenticate user** : SQL query to retrieve user by given username and password. The query must return only single arbitary column, it must have a where clause with two placeholder ::username and ::password.

- **Wrap desired endpoints to protect** : You should wrap the endpoints you want to protect with ***gologin.LoginRequired*** function in the main function.( see the example)

See the example :

```Go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/birddevelper/gologin"
	_ "github.com/go-sql-driver/mysql"
)

// static assets like CSS and JavaScript
func public() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
}

// a page in our application
func securedPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi! Welcome to secured page.")
	})
}

func main() {
	// create connection to database
	db, err := sql.Open("mysql", "root:12345@tcp(127.0.0.1:6666)/mydb")
	if err != nil {
		log.Fatal(err)
	}

	// Gologin configuration
	gologin.Configure().
		SetLoginPage("./template/login.html"). // set login page html template path
		SetSessionTimeout(90).                 // set session expiration time in seconds
		SetLoginPath("/login").                // set login http path
		// set database connection and sql query
		AuthenticateBySqlQuery(
            db,
            "select username from users where username = ::username and password = ::password")

	// instantiate http server
	mux := http.NewServeMux()

	mux.Handle("/static/", public())

	// use Gologin login handler for /login endpoint
	mux.Handle("/login", gologin.LoginHandler())

	// use Gologin logout handler for /logout endpoint
	mux.Handle("/logout", gologin.LogoutHandler())

	// the pages/endpoints that we need to protect should be wrapped with gologin.LoginRequired
	mux.Handle("/mySecuredPage", gologin.LoginRequired(securedPage()))

	// server configuration
	addr := ":8080"
	server := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// start listening to network
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("main: couldn't start simple server: %v\n", err)
	}
}

```

It is mandatory to name the login form's username input as "username" and password input as "password". Note that the form must send form data as post to the same url (set no action attribute).

Html template for login page :

```HTML
{{define "login"}}
<html>
    <body>
        <H2>
            Login Page
        </H2>
        <form method="post">
            <!-- username input with "username" name -->
            <input type="text" name="username" />
            <input type="password" name="password" />
            <input type="submit" value="Login" />
        </form>
    </body>
</html>

{{end}}

```

You can also store data in in-memory session storage in any type during user's session with **SetSession** function, and retrieve it back by **GetSession** function.

```Go

// get the session data, the request parameter is *http.Request
age, err : = gologin.GetSession("age", request)

// as the GetSession returns type is interface{}, we should specify the exact type of the session entry
fmt.Printf("Your age is " + age.(int))

```

**GetDataReturnedByAuthQuery** function returns the data of the column you specified in authentication SQL query. And with **GetCurrentUsername** you can get the current user's username.

```Go

// get the current user's username, the request parameter is *http.Request
username, err : = gologin.GetCurrentUsername(request)

fmt.Printf("Welcome " + username)

```

## Todo list

- Implement role managemnet and authorization
- mongoDB support