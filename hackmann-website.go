package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
)

const registrationsFilePath = "registrations.csv"

var staticFilePaths = []string{
	"style/style.css",
	"script/script.js",
}

var csvWriter *csv.Writer

func serveSingle(pattern string, filePath string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filePath)
	})
}

func main() {
	registrationsFile, err := os.OpenFile(registrationsFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer registrationsFile.Close()
	csvWriter = csv.NewWriter(registrationsFile)

	for _, filePath := range staticFilePaths {
		serveSingle("/"+filePath, filePath)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/register", registerHandler)
	if err := http.ListenAndServe("0.0.0.0:9000", nil); err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		http.ServeFile(w, r, "index.html")
	default:
		http.NotFound(w, r)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {

}
