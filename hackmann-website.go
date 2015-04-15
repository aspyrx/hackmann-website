package main

import (
    "bytes"
    "compress/gzip"
	"database/sql"
	_ "github.com/lib/pq"
    "io"
	"log"
	"net/http"
    "os"
    "path/filepath"
	"regexp"
    "strings"
)

type fileCache struct {
	buf bytes.Buffer
    gzipped bytes.Buffer
}

func (f *fileCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
        w.Header().Set("Content-Encoding", "gzip")
        w.Write(f.gzipped.Bytes())
    } else {
        w.Write(f.buf.Bytes())
    }
}

func NewCache(fileName string) *fileCache {
	f, err := os.Open(fileName)
    defer f.Close()
	if err != nil {
		log.Fatalf("couldn't open file: %v", err)
	}

	ret := &fileCache{}
    gzipper := gzip.NewWriter(&ret.gzipped)
    defer gzipper.Close()
	_, err = io.Copy(io.MultiWriter(&ret.buf, gzipper), f)
    if err != nil {
		log.Fatalf("couldn't read file: %v", err)
	}

	return ret
}

const staticFileDirectory = "dist"

var db *sql.DB
var staticFiles = map[string]*fileCache {}

func main() {
	var err error
	db, err = sql.Open("postgres", "connect_timeout=5 user=hackmann password='hackmann' dbname=hackmann sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
    
    staticFilePaths := []string {}
    err = filepath.Walk(staticFileDirectory, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            // deal with Windows path separator
            staticFilePaths = append(staticFilePaths, path)
        }
        return err
    })
    if err != nil {
        log.Fatal(err)
    }
    
    for _, path := range staticFilePaths {
        staticFiles[path] = NewCache(path)
    }

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/register", registerHandler)
	if err := http.ListenAndServe("0.0.0.0:80", nil); err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    var path string
    switch (r.URL.Path) {
        case "/":
            path = filepath.FromSlash(staticFileDirectory + "/index.xhtml")
        default:
            path = filepath.FromSlash(staticFileDirectory + r.URL.Path)
    }
    
    file, ok := staticFiles[path]
    if !ok {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    var contentType string
    switch (filepath.Ext(path)) {
        case ".xhtml":
            contentType = "application/xhtml+xml"
        case ".css":
            contentType = "text/css"
        case ".js":
            contentType = "application/javascript"
        case ".svg":
            contentType = "image/svg+xml"
        default:
            contentType = "text/plain"
    }
    w.Header().Set("Content-Type", contentType)
    
    file.ServeHTTP(w, r)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var firstname, lastname, school, email = r.PostForm.Get("firstname"), r.PostForm.Get("lastname"), r.PostForm.Get("school"), r.PostForm.Get("email")
	if !(len(firstname) > 0 || len(lastname) > 0 || len(school) > 0 || isValidEmail(email)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow(`SELECT id FROM registrations WHERE email = $1`, email).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		err = db.QueryRow(`INSERT INTO registrations (firstname, lastname, school, email) VALUES ($1, $2, $3, $4) RETURNING id`, firstname, lastname, school, email).Scan(&id)
	case err != nil:
		break
	default:
		err = db.QueryRow(`UPDATE registrations SET firstname = $1, lastname = $2, school = $3 WHERE id = $4 RETURNING id`, firstname, lastname, school, id).Scan(&id)
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
