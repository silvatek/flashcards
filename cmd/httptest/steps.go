package main

import (
	"flashcards/internal/test"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func addFlashcardSteps(suite *test.TestSuite) {
	suite.AddStep("the home page is accessed", whenHomePageAccessed)
	suite.AddStep("the test deck is accessed", whenTestDeckAccessed)
	suite.AddStep("the page contains `text`", thenPageContains)
	suite.AddStep("the deck has the question `question`", thenDeckHasQuestion)
}

func whenHomePageAccessed(data *test.StepData) {
	openPage("/", data)
}

func whenTestDeckAccessed(data *test.StepData) {
	openPage("/deck/TEST-CODE", data)
}

func openPage(url string, data *test.StepData) {
	response, err := http.Get(data.Suite.BaseUrl + url)

	if err == nil {
		if response.StatusCode < 400 {
			data.Suite.CurrentPageDoc, _ = goquery.NewDocumentFromReader(response.Body)
		} else {
			data.Suite.ReportError("Http request failed, %s = %d", url, response.StatusCode)
		}
	} else {
		data.Suite.ReportError("[%v]\n", err)
	}
}

func thenDeckHasQuestion(data *test.StepData) {
	if data.Suite.CurrentPageDoc == nil {
		data.Suite.ReportError("No current page document")
		return
	}

	found := false
	for _, question := range data.Suite.CurrentPageDoc.Find("span.question").Nodes {
		questionText := question.FirstChild.Data
		if strings.Contains(questionText, data.Step.Values[0]) {
			found = true
			break
		}
	}
	if !found {
		data.Suite.ReportError("Did not find question %s", data.Text)
	}
}

func thenPageContains(data *test.StepData) {
	if data.Suite.CurrentPageDoc == nil {
		data.Suite.ReportError("No current page document")
		return
	}

	pageContent := data.Suite.CurrentPageDoc.Text()

	if !strings.Contains(pageContent, data.Step.Values[0]) {
		data.Suite.ReportError("Page does not contain text: %s", data.Step.Values[0])
	}
}
