package test

import (
	"bufio"
	"fmt"
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
}

func RunScript(scanner *bufio.Scanner, suite TestSuite) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0:1] == "#" {
			continue
		}
		if strings.HasPrefix(line, "Scenario: ") {
			suite.NewScenario(strings.TrimSpace(strings.TrimPrefix(line, "Scenario: ")))
		} else if strings.HasPrefix(line, "Given ") {
			suite.DoStep("GIVEN", strings.TrimPrefix(line, "Given "))
		} else if strings.HasPrefix(line, "When ") {
			suite.DoStep("WHEN", strings.TrimPrefix(line, "When "))
		} else if strings.HasPrefix(line, "Then ") {
			suite.DoStep("THEN", strings.TrimPrefix(line, "Then "))
		} else if strings.HasPrefix(line, "And ") {
			suite.DoStep("AND", strings.TrimPrefix(line, "And "))
		} else {
			fmt.Println(line)
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

func (suite *TestSuite) NewScenario(title string) {
	suite.ReportScenarioResult()
	suite.Report("-----------")
	suite.Report("SCENARIO %s", title)
	suite.CurrentPageDoc = nil
	suite.ScenarioStatus = "NEW"
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

	for _, step := range suite.Steps {
		if step.parsed.Matches(parsedStep) {
			fmt.Printf("%8s %s\n", stepType, text)
			step.impl(&data)
			return
		}
	}
	suite.ReportError("No step definition matches: %s\n", text)
}
