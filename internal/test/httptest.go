package test

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type StepData struct {
	Text  string
	Step  ParsedStep
	Suite *TestSuite
}

type StepDefinition struct {
	Text   string
	parsed ParsedStep
	impl   func(data *StepData)
}

type TestSuite struct {
	Steps          []StepDefinition
	BaseUrl        string
	CurrentPageDoc *goquery.Document
	ScenarioStatus string
	TotalFeatures  int
	TotalScenarios int
	TotalChecks    int
	ChecksPassed   int
	IsChecking     bool
}

func RunScript(scanner *bufio.Scanner, suite *TestSuite) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0:1] == "#" {
			// no-op
		} else if strings.HasPrefix(line, "Feature: ") {
			suite.NewFeature(strings.TrimPrefix(line, "Feature: "))
		} else if strings.HasPrefix(line, "Scenario: ") {
			suite.NewScenario(strings.TrimPrefix(line, "Scenario: "))
		} else if strings.HasPrefix(line, "Given ") {
			suite.DoStep("GIVEN", strings.TrimPrefix(line, "Given "))
		} else if strings.HasPrefix(line, "When ") {
			suite.DoStep("WHEN", strings.TrimPrefix(line, "When "))
		} else if strings.HasPrefix(line, "Then ") {
			suite.IsChecking = true
			suite.DoStep("THEN", strings.TrimPrefix(line, "Then "))
		} else if strings.HasPrefix(line, "And ") {
			suite.DoStep("AND", strings.TrimPrefix(line, "And "))
		} else {
			fmt.Printf("Unsupported script element: %s\n", line)
		}
	}

	suite.ReportScenarioResult()
}

func (suite *TestSuite) AddStep(text string, impl func(*StepData)) {
	if suite.Steps == nil {
		suite.Steps = make([]StepDefinition, 0)
	}
	suite.Steps = append(suite.Steps, StepDefinition{text, ParseStep(text), impl})
}

func (suite *TestSuite) Report(text string, args ...any) {
	fmt.Printf(text+"\n", args...)
}

func (suite *TestSuite) ReportError(text string, args ...any) {
	fmt.Printf("   ERROR "+text+"\n", args...)
	suite.ScenarioStatus = "FAILED"
}

func (suite *TestSuite) ReportScenarioResult() {
	if suite.ScenarioStatus == "" {
		// Do nothing
	} else if suite.ScenarioStatus == "FAILED" {
		fmt.Println("  STATUS Scenario failed")
	} else {
		fmt.Println("  STATUS Scenario passed")
	}
}

func (suite *TestSuite) NewFeature(title string) {
	suite.TotalFeatures += 1
}

func (suite *TestSuite) NewScenario(title string) {
	suite.ReportScenarioResult()
	suite.Report("-----------")
	suite.Report("SCENARIO %s", title)
	suite.CurrentPageDoc = nil
	suite.ScenarioStatus = "NEW"
	suite.IsChecking = false
	suite.TotalScenarios += 1
}

func (suite *TestSuite) DoStep(stepType string, text string) {
	var data StepData
	data.Text = text
	data.Suite = suite
	data.Step = ParseStep(text)

	if data.Suite.ScenarioStatus == "FAILED" {
		return
	}

	parsedStep := ParseStep(text)

	if suite.IsChecking {
		suite.TotalChecks += 1
	}

	for _, step := range suite.Steps {
		if step.parsed.Matches(parsedStep) {
			fmt.Printf("%8s %s\n", stepType, text)
			step.impl(&data)
			return
		}
	}
	suite.ReportError("No step definition matches: %s\n", text)
}

func (suite *TestSuite) OpenPage(path string) {
	response, err := http.Get(suite.BaseUrl + path)

	if err == nil {
		if response.StatusCode < 400 {
			suite.CurrentPageDoc, _ = goquery.NewDocumentFromReader(response.Body)
		} else {
			suite.ReportError("Http request failed, %s = %d", path, response.StatusCode)
		}
	} else {
		suite.ReportError("[%v]\n", err)
	}
}

func (suite *TestSuite) Summary() {
	fmt.Println("===========")
	fmt.Printf("Features tested: %d\n", suite.TotalFeatures)
	fmt.Printf("Test scenarios:  %d\n", suite.TotalScenarios)
	passrate := 100 * suite.ChecksPassed / suite.TotalChecks
	fmt.Printf("Checks passed:   %d out of %d (%d%%)\n", suite.ChecksPassed, suite.TotalChecks, passrate)
	fmt.Println("===========")
}
