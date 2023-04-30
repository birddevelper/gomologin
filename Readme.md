# What is Gomologin

<p align="center">
<img src="https://mshaeri.com/blog/wp-content/uploads/2022/04/gologin.png"  height="200" >
</p>

**Gomologin** is an easy to setup professional login manager for Go web applications. It helps you protect your application resources from unattended, unauthenticated or unauthorized access. Currently it works with SQL databases authentication. It is flexible, you can use it with any user/roles table structure in database.

## How to setup

Get the package with following command :

```bash
go get github.com/birddevelper/gomologin

```

## How to use

(You can read detailed tutorial with code example at [gomoloign user guide](https://mshaeri.com/blog/golang-login-manager-with-gomologin-package/))

You can easily setup and customize login process with **configure()** function. You should specify following paramters to make the Gomologin ready to start:

- **Login page** : path to html template. Default path is ***./template/login.html***, note that the template must be defined as ****"login"**** with ***{{define "login"}}*** at the begining line

- **Login path** : login http path. Default path is ***/login***

- **Session timeout** : Number of seconds before the session expires. Default value is 120 seconds.

- **Password encryption** : Password encryption function to apply on password before it compare with password stored in db. Default is ***EncNoEncrypt***

- **SQL connection, and SQL query to authenticate user and fetch roles** : 2 SQL queries to retrieve user and its roles by given username and password. The authentication query must return only single arbitary column, it must have a where clause with two placeholder ::username and ::password. And the query for retrieving user's roles must return only the text column of role name.

- **Wrap desired endpoints to protect** : You should wrap the endpoints you want to protect with ***gomologin.LoginRequired*** or ***gomologin.RolesRequired*** function in the main function.( see the example)

***gomologin.LoginRequired*** requires user to be authenticated for accessing the wrapped endpoint/page.

***gomologin.RolesRequired*** requires user to have specified roles in addition to be authenticated.

See the example :

```Go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/birddevelper/gomologin"
	_ "github.com/go-sql-driver/mysql"
)

// static assets like CSS and JavaScript
func public() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
}

// a page in our application, it needs user only be authenticated
func securedPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi! Welcome to secured page.")
	})
}

// another page in our application, it needs user be authenticated and have ADMIN role
func securedPage2() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi! Welcome to very secured page.")
	})
}


func main() {
	// create connection to database
	db, err := sql.Open("mysql", "root:12345@tcp(127.0.0.1:6666)/mydb")
	if err != nil {
		log.Fatal(err)
	}

	// Gomologin configuration
	gomologin.Configure().
		SetLoginPage("./template/login.html"). // set login page html template path
		SetSessionTimeout(90).                 // set session expiration time in seconds
		SetLoginPath("/login").                // set login http path
		// set database connection and sql query
		AuthenticateBySqlQuery(
			db,
			"select id from users where username = ::username and password = ::password", // authentication query
			"select role from user_roles where userid = (select id from users where username = ::username)") // fetch user's roles

	// instantiate http server
	mux := http.NewServeMux()

	mux.Handle("/static/", public())

	// use Gomologin login handler for /login endpoint
	mux.Handle("/login", gomologin.LoginHandler())

	// the pages/endpoints that we need to protect should be wrapped with gomologin.LoginRequired
	mux.Handle("/mySecuredPage", gomologin.LoginRequired(securedPage()))

	mux.Handle("/mySecuredPage2", gomologin.RolesRequired(securedPage2(),"ADMIN"))

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

It is mandatory to set the login form's username input as "username" and password input as "password". Note that the form must send form data as post to the same url (set no action attribute).

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
func securedPage2() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the session data, the request parameter is *http.Request
		age, err : = gomologin.GetSession("age", request)

		// as the GetSession returns type is interface{}, we should specify the exact type of the session entry
		fmt.Printf("Your age is " + age.(int))
	})
}
```

**GetDataReturnedByAuthQuery** function returns the data of the column you specified in authentication SQL query. And with **GetCurrentUsername** you can get the current user's username.

```Go
func securedPage2() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get the current user's username, the request parameter is *http.Request
			username : = gomologin.GetCurrentUsername(request)

			fmt.Printf("Welcome " + username)
	})
}
```

To logout users direct them to your **login url + ?logout=yes** for example if your login url is **/login** your application logout url will be **/login?logout=yes**

You can read detailed tutorial with code example at [gomoloign user guide](https://mshaeri.com/blog/golang-login-manager-with-gomologin-package/)


## Todo list

- mongoDB support
