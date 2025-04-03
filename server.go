package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	ContentType = "Content-Type"
)

// Load style.css from the embedded file system
//
//go:embed www/style.css
var styleCSS []byte

// Embed the HTML template for the main page
//
//go:embed www/template.html
var htmlTemplate string

// Embed the font
//
//go:embed www/Yantramanav-Regular.ttf
var font []byte

// serve starts the HTTP server and handles requests
func serve(portNumber string, interval time.Duration, logRequests bool) {
	// Parse the HTML template
	tmp, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(ContentType, "text/css")
		w.Write(styleCSS)
	})

	http.HandleFunc("/Yantramanav-Regular.ttf", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(ContentType, "font/ttf")
		w.Write(font)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if logRequests {
			log.Println(r.RemoteAddr, r.URL.Path)
		}
		w.Header().Set(ContentType, "text/html; charset=utf-8")
		// Create a struct to pass both the launches and the time remaining to the template
		data := struct {
			Launches      []*Launch
			TimeRemaining string
		}{
			Launches:      launchesCache,
			TimeRemaining: timeUntilNextRefresh(interval),
		}

		if err := tmp.Execute(w, data); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	})

	log.Println("Starting server on " + portNumber)
	log.Fatal(http.ListenAndServe(":"+portNumber, nil))
}
