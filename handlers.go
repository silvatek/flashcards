package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/mux"
	"github.com/skip2/go-qrcode"
)

type pageData struct {
	Message    string
	Error      string
	Deck       Deck
	Card       Card
	Show       string
	Share      string
	Question   template.HTML
	Answer     template.HTML
	Hint       template.HTML
	FormAction string
	History    []string
}

const HISTORY_COOKIE = "deckHistory"

func applicationRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", homePage)
	r.HandleFunc("/decks", deckRedirect)
	r.HandleFunc("/deck/{id}/card/{card}", cardPage)
	r.HandleFunc("/deck/{id}", deckPage)
	r.HandleFunc("/random", randomCard)
	r.HandleFunc("/newcard", addCard)
	r.HandleFunc("/editcard", editCard)
	r.HandleFunc("/newdeck", newDeck)
	r.HandleFunc("/error", errorPage)
	r.HandleFunc("/qrcode", qrCodeGenerator)

	addStaticAssetRouter(r)

	return r
}

func addStaticAssetRouter(r *mux.Router) {
	fs := http.FileServer(http.Dir("template/static"))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
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
		History: getHistory(HISTORY_COOKIE, r).entries,
	}

	showTemplatePage("index", data, w)
}

func deckRedirect(w http.ResponseWriter, r *http.Request) {
	deckId := r.FormValue("deck")
	deckUrl := fmt.Sprintf("/deck/%s", strings.ToUpper(deckId))
	logs.debug("Redirecting to %s", deckUrl)
	http.Redirect(w, r, deckUrl, http.StatusSeeOther)
}

func deckPage(w http.ResponseWriter, r *http.Request) {
	logs.debug("Deck page %s", r.RequestURI)
	deckID := mux.Vars(r)["id"]

	logs.debug("Showing deck %s", deckID)

	var shareUrl string
	if r.FormValue("share") == "true" {
		shareUrl = deckUrl(r, deckID)
	}

	data := pageData{
		Deck:  dataStore.getDeck(context.Background(), deckID),
		Share: shareUrl,
	}

	if data.Deck.ID != deckID {
		http.Redirect(w, r, "/error?code=2001", http.StatusSeeOther)
		return
	}

	history := getHistory(HISTORY_COOKIE, r)
	history.push(deckID)
	history.setCookie(w)

	showTemplatePage("deck", data, w)

}

func deckUrl(r *http.Request, deckID string) string {
	if r.Host == "localhost:8080" {
		return fmt.Sprintf("http://localhost:8080/deck/%s", deckID)
	} else {
		return fmt.Sprintf("https://%s/deck/%s", r.Host, deckID)
	}
}

func cardPage(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["id"]
	cardID := mux.Vars(r)["card"]

	logs.debug("Showing card %s from deck %s", cardID, deckID)

	deck := dataStore.getDeck(context.Background(), deckID)

	if deck.ID != deckID {
		http.Redirect(w, r, "/error?code=2001", http.StatusSeeOther)
		return
	}

	card := deck.getCard(cardID)

	if card.ID != cardID {
		http.Redirect(w, r, "/error?code=2002", http.StatusSeeOther)
		return
	}

	show := strings.ToLower(r.FormValue("answer"))
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
	deckId := strings.ToUpper(r.FormValue("deck"))

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
		deckID := strings.ToUpper(r.FormValue("deck"))
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
		deckID := strings.ToUpper(r.FormValue("deck"))
		cardID := strings.ToUpper(r.FormValue("card"))
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

	deck := Deck{
		ID:    randomDeckId(),
		Title: r.Form.Get("title"),
	}

	if !dataStore.isValidAuthor(r.Form.Get("author")) {
		http.Redirect(w, r, "/error?code=3001", http.StatusSeeOther)
		return
	}

	logs.info("Creating deck %s with title %s", deck.ID, deck.Title)

	dataStore.putDeck(context.Background(), deck.ID, deck)

	http.Redirect(w, r, "/deck/"+deck.ID, http.StatusSeeOther)
}

func errorPage(w http.ResponseWriter, r *http.Request) {
	errorCode := r.FormValue("code")
	data := pageData{
		Error: errorText(errorCode),
	}
	logs.info("Showing error page %s %s", errorCode, data.Error)
	showTemplatePage("error", data, w)
}

func qrCodeGenerator(w http.ResponseWriter, r *http.Request) {
	deckID := strings.ToUpper(r.FormValue("deck"))

	gameUrl := deckUrl(r, deckID)

	tempFileName := os.Getenv("TMPDIR") + "/" + deckID + ".png"

	qrcode.WriteFile(gameUrl, qrcode.High, 320, tempFileName)
	defer os.Remove(tempFileName)

	content, _ := os.ReadFile(tempFileName)

	headers := w.Header()
	headers.Add("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}
