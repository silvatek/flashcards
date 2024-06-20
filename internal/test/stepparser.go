package test

import (
	"strings"
	"text/scanner"
)

type ParsedStep struct {
	Unparsed   string
	Comparable string
	Values     []string
}

func ParseStep(text string) ParsedStep {
	var step ParsedStep
	step.Unparsed = text
	step.Values = make([]string, 0)

	var s scanner.Scanner
	s.Init(strings.NewReader(text))

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		token := s.TokenText()
		if token[0] == '`' {
			step.Values = append(step.Values, token[1:len(token)-1])
			token = "*"
		}
		step.Comparable = step.Comparable + " " + strings.TrimSpace(token)
	}
	step.Comparable = strings.TrimSpace(step.Comparable)
	return step
}

func (s *ParsedStep) Matches(s2 ParsedStep) bool {
	return s.Comparable == s2.Comparable
}
