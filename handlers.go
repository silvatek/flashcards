package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func addHandlers() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/decks/", deckPage)

	addStaticAssetHandler()
}

func addStaticAssetHandler() {
	fs := http.FileServer(http.Dir("template/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
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

// show home/index page
func homePage(w http.ResponseWriter, r *http.Request) {
	logs.info("Received request: %s %s", r.Method, r.URL.Path)

	data := pageData{
		Message: "Ice Hockey Scoresheet",
	}

	showTemplatePage("index", data, w)
}

func deckPage(w http.ResponseWriter, r *http.Request) {

	deckID := lastPathElement(r.RequestURI)

	logs.debug("Showing deck %s", deckID)

	data := pageData{
		Deck: dataStore.getDeck(context.Background(), deckID),
	}

	showTemplatePage("deck", data, w)
}

func lastPathElement(uri string) string {
	// strip query parameters
	queryStart := strings.Index(uri, "?")
	if queryStart > -1 {
		uri = uri[:queryStart]
	}
	// return everything after the last slash
	lastSlash := strings.LastIndex(uri, "/")
	if lastSlash == -1 {
		return uri
	}
	return uri[lastSlash+1:]
}
