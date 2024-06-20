package main

import (
	"flashcards/internal/test"
	"strings"
)

func addFlashcardSteps(suite *test.TestSuite) {
	suite.AddStep("the home page is accessed", whenHomePageAccessed)
	suite.AddStep("the test deck is accessed", whenTestDeckAccessed)
	suite.AddStep("the page contains `text`", thenPageContains)
	suite.AddStep("the deck has the question `question`", thenDeckHasQuestion)
}

func whenHomePageAccessed(data *test.StepData) {
	data.Suite.OpenPage("/")
}

func whenTestDeckAccessed(data *test.StepData) {
	data.Suite.OpenPage("/deck/TEST-CODE")
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
			data.Suite.ChecksPassed += 1
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

	if strings.Contains(pageContent, data.Step.Values[0]) {
		data.Suite.ChecksPassed += 1
	} else {
		data.Suite.ReportError("Page does not contain text: %s", data.Step.Values[0])
	}
}
