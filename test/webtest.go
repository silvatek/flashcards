package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

type WebTest struct {
	t        *testing.T
	Response *httptest.ResponseRecorder
	Request  *http.Request
	doc      *goquery.Document
	path     string
	success  bool
}

func NewWebTest(t *testing.T) WebTest {
	return WebTest{
		t:        t,
		Response: httptest.NewRecorder(),
		success:  true,
	}
}

func (wt *WebTest) SendGet(path string) {
	wt.path = path
	wt.Request = httptest.NewRequest(http.MethodGet, path, nil)
}

func (wt *WebTest) ShowBodyOnFail() {
	if !wt.success {
		wt.t.Log(wt.doc.Html())
	}
}

func (wt *WebTest) AssertSuccess() {
	if wt.Response.Code >= 400 {
		wt.success = false
		wt.t.Errorf("Non-success response code (%d) for path %s", wt.Response.Code, wt.path)
	}
}

func (wt *WebTest) AssertRedirectTo(expectedTarget string) {
	if wt.Response.Code != http.StatusSeeOther {
		wt.t.Errorf("Non-redirect response code (%d) for path %s", wt.Response.Code, wt.path)
		return
	}
	redirects := wt.Response.Header().Values("Location")
	if len(redirects) == 0 {
		wt.t.Errorf("No redirect header for path %s", wt.path)
		return
	}
	redirectTo := redirects[0]
	if redirectTo != expectedTarget {
		wt.t.Errorf("Unexpected redirect target for path %s, %s != %s", wt.path, redirectTo, expectedTarget)
	}
}

func (wt *WebTest) AssertBodyContains(query string, expected string) {
	if wt.doc == nil {
		wt.doc, _ = goquery.NewDocumentFromReader(wt.Response.Body)
	}
	text := wt.doc.Find(query).Text()
	if !strings.Contains(text, expected) {
		wt.success = false
		wt.t.Errorf("Did not find %s in %s with query %s", expected, wt.path, query)
	}
}
