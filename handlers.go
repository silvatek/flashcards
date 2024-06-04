package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type pageData struct {
	Message    string
	Error      string
	Deck       Deck
	Card       Card
	ShowAnswer bool
}

func addHandlers() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/deck/", deckPage)
	http.HandleFunc("/random", randomCard)
	http.HandleFunc("/newcard", addCard)

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
	logs.debug("Deck page %s", r.RequestURI)
	if strings.Contains(r.RequestURI, "/card/") {
		cardPage(w, r)
	} else {
		deckID := lastPathElement(r.RequestURI)

		logs.debug("Showing deck %s", deckID)

		data := pageData{
			Deck: dataStore.getDeck(context.Background(), deckID),
		}

		showTemplatePage("deck", data, w)
	}
}

func cardPage(w http.ResponseWriter, r *http.Request) {
	cardID := lastPathElement(r.RequestURI)

	path := strings.Replace(r.RequestURI, "/card/"+cardID, "", 1)

	deckID := lastPathElement(path)

	logs.debug("Showing card %s from deck %s", cardID, deckID)

	deck := dataStore.getDeck(context.Background(), deckID)
	card := deck.getCard(cardID)

	data := pageData{
		Deck:       deck,
		Card:       card,
		ShowAnswer: queryParam(r.RequestURI, "answer") == "show",
	}
	showTemplatePage("card", data, w)
}

func randomCard(w http.ResponseWriter, r *http.Request) {
	deckId := queryParam(r.RequestURI, "deck")

	deck := dataStore.getDeck(context.Background(), deckId)

	if deck.ID == "" {
		logs.error("Could not fetch deck %s", deckId)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	logs.info("Showing random card for %s", deck.Title)

	card := deck.randomCard()

	http.Redirect(w, r, "/deck/"+deckId+"/card/"+card.ID+"?answer=hide", http.StatusSeeOther)
}

func addCard(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		deckID := r.Form.Get("deck_id")

		deck := dataStore.getDeck(context.Background(), deckID)

		card := Card{
			ID:       randomCardId(),
			DeckID:   deckID,
			Question: r.Form.Get("question"),
			Answer:   r.Form.Get("answer"),
		}

		deck.addCard(card)

		dataStore.putDeck(context.Background(), deck.ID, deck)

		http.Redirect(w, r, "/deck/"+deckID, http.StatusSeeOther)
	} else {
		deckID := queryParam(r.RequestURI, "deck")
		logs.debug("Showing new card page for %s", deckID)
		deck := dataStore.getDeck(context.Background(), deckID)
		data := pageData{
			Deck: deck,
		}
		showTemplatePage("newcard", data, w)
	}
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

func queryParam(uri string, param string) string {
	queryStart := strings.Index(uri, "?")
	if queryStart == -1 {
		return ""
	}
	uri = uri[queryStart+1:]

	paramStart := strings.Index(uri, param+"=")
	if paramStart == -1 {
		return ""
	}
	paramVal := uri[paramStart:]

	valueStart := strings.Index(uri, "=")
	paramVal = paramVal[valueStart+1:]

	nextStart := strings.Index(paramVal, "&")
	if nextStart > 0 {
		paramVal = paramVal[0:nextStart]
	}

	return paramVal
}
