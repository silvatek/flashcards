package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

type WebTest struct {
	t       *testing.T
	w       *httptest.ResponseRecorder
	r       *http.Request
	doc     *goquery.Document
	path    string
	success bool
}

func webTest(t *testing.T) WebTest {
	return WebTest{
		t:       t,
		w:       httptest.NewRecorder(),
		success: true,
	}
}

func (wt *WebTest) sendGet(path string) {
	wt.path = path
	wt.r = httptest.NewRequest(http.MethodGet, path, nil)
}

func (wt *WebTest) showBodyOnFail() {
	if !wt.success {
		wt.t.Log(wt.doc.Html())
	}
}

func (wt *WebTest) assertSuccess() {
	if wt.w.Code >= 400 {
		wt.success = false
		wt.t.Errorf("Non-success response code (%d) for path %s", wt.w.Code, wt.path)
	}
}

func (wt *WebTest) assertRedirectTo(expectedTarget string) {
	if wt.w.Code != http.StatusSeeOther {
		wt.t.Errorf("Non-redirect response code (%d) for path %s", wt.w.Code, wt.path)
		return
	}
	redirects := wt.w.Header().Values("Location")
	if len(redirects) == 0 {
		wt.t.Errorf("No redirect header for path %s", wt.path)
		return
	}
	redirectTo := redirects[0]
	if redirectTo != expectedTarget {
		wt.t.Errorf("Unexpected redirect target for path %s, %s != %s", wt.path, redirectTo, expectedTarget)
	}
}

func (wt *WebTest) assertBodyContains(query string, expected string) {
	if wt.doc == nil {
		wt.doc, _ = goquery.NewDocumentFromReader(wt.w.Body)
	}
	text := wt.doc.Find(query).Text()
	if !strings.Contains(text, expected) {
		wt.success = false
		wt.t.Errorf("Did not find %s in %s with query %s", expected, wt.path, query)
	}
}
