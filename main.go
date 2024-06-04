package main

import (
	"html/template"
	"net/http"
	"os"
)

const defaultAddr = ":8080"

type pageData struct {
	Message string
	Error   string
}

var logs Logger

// main starts an http server on the $PORT environment variable.
func main() {
	logs.init()

	addr := defaultAddr
	// $PORT environment variable is provided in the Kubernetes deployment.
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	logs.info("Server listening on port %s", addr)

	addHandlers()

	if err := http.ListenAndServe(addr, nil); err != nil {
		logs.error("Server listening error: %+v", err)
		os.Exit(-5)
	}
}

func addHandlers() {
	http.HandleFunc("/", homePage)

	addStaticAssetHandler()
}

func addStaticAssetHandler() {
	fs := http.FileServer(http.Dir("template/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

// show home/index page
func homePage(w http.ResponseWriter, r *http.Request) {
	logs.info("Received request: %s %s", r.Method, r.URL.Path)

	data := pageData{
		Message: "Ice Hockey Scoresheet",
	}

	showTemplatePage("index", data, w)
}

func showTemplatePage(templateName string, data any, w http.ResponseWriter) {
	t, err := template.ParseFiles("template/" + templateName + ".html")
	if err != nil {
		logs.error("Error parsing template: %+v", err)
		os.Exit(-2)
	}

	if err := t.Execute(w, data); err != nil {
		msg := http.StatusText(http.StatusInternalServerError)
		logs.debug("template.Execute: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}
}
