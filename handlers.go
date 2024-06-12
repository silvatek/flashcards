package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type pageData struct {
	Message    string
	Error      string
	Deck       Deck
	Card       Card
	Show       string
	Question   template.HTML
	Answer     template.HTML
	Hint       template.HTML
	FormAction string
}

func addHandlers() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/decks", deckRedirect)
	http.HandleFunc("/deck/", deckPage)
	http.HandleFunc("/random", randomCard)
	http.HandleFunc("/newcard", addCard)
	http.HandleFunc("/editcard", editCard)
	http.HandleFunc("/newdeck", newDeck)

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
	logs.debug("Received request: %s %s", r.Method, r.URL.Path)

	data := pageData{
		Message: "Fashcards",
	}

	showTemplatePage("index", data, w)
}

func deckRedirect(w http.ResponseWriter, r *http.Request) {
	deckId := strings.ToUpper(queryParam(r.RequestURI, "deck"))
	http.Redirect(w, r, "/deck/"+deckId, http.StatusSeeOther)
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

	show := strings.ToLower(queryParam(r.RequestURI, "answer"))
	if show == "" {
		show = "none"
	}

	data := pageData{
		Deck:     deck,
		Card:     card,
		Show:     show,
		Question: renderMarkdown(card.Question),
		Answer:   renderMarkdown(card.Answer),
		Hint:     renderMarkdown(card.Hint),
	}
	showTemplatePage("card", data, w)
}

func renderMarkdown(source string) template.HTML {
	// create markdown parser with extensions
	extensions := parser.Tables | parser.Strikethrough
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(source))

	// create HTML renderer with extensions
	htmlFlags := html.SkipLinks | html.SkipImages | html.SkipHTML
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	rendered := string(markdown.Render(doc, renderer)[:])

	// convert the HTML into a template fragment
	htmlFragment := template.HTML(rendered)
	return htmlFragment
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
			Hint:     r.Form.Get("hint"),
		}

		deck.addCard(card)

		dataStore.putDeck(context.Background(), deck.ID, deck)

		http.Redirect(w, r, "/deck/"+deckID, http.StatusSeeOther)
	} else {
		deckID := queryParam(r.RequestURI, "deck")
		logs.debug("Showing new card page for %s", deckID)
		deck := dataStore.getDeck(context.Background(), deckID)
		data := pageData{
			Deck:       deck,
			Card:       *new(Card),
			FormAction: "/newcard",
		}
		showTemplatePage("editcard", data, w)
	}
}

func editCard(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		deckID := r.Form.Get("deck_id")
		cardID := r.Form.Get("card_id")

		logs.info("Received edit for card %s in deck %s", cardID, deckID)

		deck := dataStore.getDeck(context.Background(), deckID)

		card := deck.getCard(cardID)

		card.Question = r.Form.Get("question")
		card.Answer = r.Form.Get("answer")
		card.Hint = r.Form.Get("hint")

		logs.debug("New answer: %s", card.Answer)

		deck.putCard(card.ID, card)

		dataStore.putDeck(context.Background(), deck.ID, deck)

		http.Redirect(w, r, "/deck/"+deckID+"/card/"+cardID+"?answer=show", http.StatusSeeOther)
	} else {
		deckID := queryParam(r.RequestURI, "deck")
		cardID := queryParam(r.RequestURI, "card")
		logs.debug("Showing edit card page for %s / %s", deckID, cardID)
		deck := dataStore.getDeck(context.Background(), deckID)
		data := pageData{
			Deck:       deck,
			Card:       deck.getCard(cardID),
			FormAction: "/editcard",
		}
		showTemplatePage("editcard", data, w)
	}
}

func newDeck(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Form.Get("author") != "guessme" {
		logs.info("Attempt to create a new deck without a valid author code")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	deck := Deck{
		ID:    randomDeckId(),
		Title: r.Form.Get("title"),
	}

	logs.info("Creating deck %s with title %s", deck.ID, deck.Title)

	dataStore.putDeck(context.Background(), deck.ID, deck)

	http.Redirect(w, r, "/deck/"+deck.ID, http.StatusSeeOther)
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
