package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"regexp"
    "strings"
)

var Db *sql.DB

func main() {
	var err error
	Db, err = sql.Open("postgres", "connect_timeout=5 user=hackmann password='hackmann' dbname=hackmann sslmode=disable")
	defer Db.Close()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/register", registerHandler)
	if err := http.ListenAndServe("0.0.0.0:80", nil); err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    const rootDirectory = "./dist"
    switch (r.URL.Path) {
        case "/":
            w.Header().Set("Content-Type", "application/xhtml+xml")
            http.ServeFile(w, r, rootDirectory + "/index.xhtml")
        default:
            if strings.HasSuffix(r.URL.Path, ".svg") {
                w.Header().Set("Content-Type", "image/svg+xml")
            }
            http.ServeFile(w, r, rootDirectory + r.URL.Path)
    }
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var firstname, lastname, school, email = r.PostForm.Get("firstname"), r.PostForm.Get("lastname"), r.PostForm.Get("school"), r.PostForm.Get("email")
	if !(len(firstname) > 0 || len(lastname) > 0 || len(school) > 0 || isValidEmail(email)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var id int
	err := Db.QueryRow(`SELECT id FROM registrations WHERE email = $1`, email).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		err = Db.QueryRow(`INSERT INTO registrations (firstname, lastname, school, email) VALUES ($1, $2, $3, $4) RETURNING id`, firstname, lastname, school, email).Scan(&id)
	case err != nil:
		break
	default:
		err = Db.QueryRow(`UPDATE registrations SET firstname = $1, lastname = $2, school = $3 WHERE id = $4 RETURNING id`, firstname, lastname, school, id).Scan(&id)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// regexp source: https://github.com/asaskevich/govalidator
func isValidEmail(email string) bool {
	return regexp.MustCompile("(?i)^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$").MatchString(email)
}
