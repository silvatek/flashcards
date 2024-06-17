# Flashcards application

For creating decks of flashcards, each with a question and answer, and practicing answering the questions.

## Implementation

Build in the Go language, initially targetted at the Google Cloud Run hosting platorm.

## Commands

`go test ./... -coverprofile=cover.out`

`go tool cover -html=cover.out`  