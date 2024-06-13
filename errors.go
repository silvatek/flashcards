package main

import "fmt"

var errors = map[string]string{
	"1001": "Unknown error",
	"2001": "Deck not found",
	"2002": "Card not found",
	"3001": "Not authorised to create new decks",
}

func errorText(errorCode string) string {
	text, ok := errors[errorCode]
	if ok {
		return text
	} else {
		return fmt.Sprintf("Unknown error %s", errorCode)
	}
}
