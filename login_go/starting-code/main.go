package main

import (
	"github.com/satori/go.uuid"
	"html/template"
	"net/http"
)

type user struct {
	UserName string
	Password string
	First    string
	Last     string
}

var tpl *template.Template
var dbUsers = map[string]user{}      // user ID, user
var dbSessions = map[string]string{} // session ID, user ID

func init() {
	tpl = template.Must(template.ParseGlob("/Users/liuzh/OneDrive/Desktop/login_go/starting-code/templates/*"))
	dbUsers["zliu112@stevens.edu"] = user{"zliu112@stevens.edu", "lzmy0309", "Jeff", "Liu"} // populate our db
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/login", login)
	http.HandleFunc("/readcookie", ReadCookieServer)
	http.HandleFunc("/wrongInformation", wrongInformation)
	http.HandleFunc("/wrongBar", wrongBar)
	http.HandleFunc("/hasCookie", hasCookie)
	http.HandleFunc("/cookieEmpty", cookieEmpty)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func ReadCookieServer(w http.ResponseWriter, req *http.Request) {
	c, _ := req.Cookie("session")
	un, ok := dbSessions[c.Value]
	if !ok {
		http.Redirect(w, req, "/cookieEmpty", http.StatusSeeOther)
		return
	}
	u := dbUsers[un]
	tpl.ExecuteTemplate(w, "hasCookie.gohtml", u)
}


func index(w http.ResponseWriter, req *http.Request) {

	// get cookie
	c, err := req.Cookie("session")
	if err != nil {
		sID, _ := uuid.NewV4() // individual Unique Identifier
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
	}


	var u user
	if un, ok := dbSessions[c.Value]; ok {
		u = dbUsers[un]
	}


	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")
		u = user{un, p, f, l}
		dbSessions[c.Value] = un
		dbUsers[un] = u
	}

	tpl.ExecuteTemplate(w, "index.gohtml", u)
}

func cookieEmpty(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "cookieEmpty.gohtml", nil)
}

func hasCookie(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w, "hasCookie.gohtml", nil)
}
func wrongInformation(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "wrongInformation.gohtml", nil)
}

func wrongBar(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "wrongBar.gohtml", nil)
}

func bar(w http.ResponseWriter, req *http.Request) {

	// get cookie
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/wrongBar", http.StatusSeeOther)
		return
	}
	un, ok := dbSessions[c.Value]
	if !ok {
		http.Redirect(w, req, "/wrongBar", http.StatusSeeOther)
		return
	}

	u := dbUsers[un]
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}

func login(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		// does this username exist
		u, ok := dbUsers[un]
		if !ok {
			//http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, req, "/wrongInformation", http.StatusSeeOther)
			return
		}
		// does the username/password combo have a match?
		if u.Password != p {
			http.Redirect(w, req, "/wrongInformation", http.StatusSeeOther)
			return
		}

		//create a session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)


}
