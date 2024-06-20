package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"flashcards/internal/test"
)

func main() {
	fmt.Println("Post-deployment Flashcard Tests")

	var suite test.TestSuite
	addFlashcardSteps(&suite)

	flag.StringVar(&suite.BaseUrl, "h", "http://localhost:8080", "Base URL, no trailing slash")
	flag.Parse()

	featureFile := flag.Arg(0)

	if featureFile == "" {
		fmt.Println("A test script (feature file) must be specified on the command line.")
		os.Exit(-1)
	}

	f, _ := os.Open(featureFile)
	defer f.Close()

	test.RunScript(bufio.NewScanner(f), suite)
}
