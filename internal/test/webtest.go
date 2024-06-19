package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type WebTest struct {
	t        *testing.T
	Response *httptest.ResponseRecorder
	Request  *http.Request
	doc      *goquery.Document
	path     string
	method   string
	success  bool
	router   mux.Router
}

func NewWebTest(t *testing.T, router mux.Router) WebTest {
	return WebTest{
		t:        t,
		Response: httptest.NewRecorder(),
		success:  true,
		router:   router,
	}
}

func (wt *WebTest) SendGet(path string) {
	wt.method = http.MethodGet
	wt.path = path
	wt.Request = httptest.NewRequest(wt.method, wt.path, nil)
	wt.router.ServeHTTP(wt.Response, wt.Request)
}

func (wt *WebTest) SendPost(path string, fields map[string]string) {
	wt.method = http.MethodPost
	wt.path = path
	body := formPostBody(fields)
	wt.Request = httptest.NewRequest(wt.method, wt.path, &body)
	wt.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	wt.router.ServeHTTP(wt.Response, wt.Request)
}

func formPostBody(fields map[string]string) bytes.Buffer {
	var buf bytes.Buffer
	empty := true
	for key, value := range fields {
		if !empty {
			fmt.Fprint(&buf, "&")
		}
		fmt.Fprintf(&buf, "%s=%s", key, value)
		empty = false
	}
	return buf
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
		wt.success = false
		wt.t.Errorf("Non-redirect response code (%d) for path %s", wt.Response.Code, wt.path)
		return
	}
	redirectTo := wt.RedirectTarget()
	if redirectTo != expectedTarget {
		wt.success = false
		wt.t.Errorf("Unexpected redirect target for path %s, %s != %s", wt.path, redirectTo, expectedTarget)
	}
}

func (wt *WebTest) AssertRedirectToPrefix(expectedTargetPrefix string) {
	if wt.Response.Code != http.StatusSeeOther {
		wt.success = false
		wt.t.Errorf("Non-redirect response code (%d) for path %s", wt.Response.Code, wt.path)
		return
	}
	redirectTo := wt.RedirectTarget()
	if !strings.HasPrefix(redirectTo, expectedTargetPrefix) {
		wt.success = false
		wt.t.Errorf("Unexpected redirect target prefix for path %s, %s != %s", wt.path, redirectTo, expectedTargetPrefix)
	}
}

func (wt *WebTest) RedirectTarget() string {
	redirects := wt.Response.Header().Values("Location")
	if len(redirects) == 0 {
		wt.success = false
		wt.t.Errorf("No redirect header for path %s", wt.path)
		return ""
	}
	return redirects[0]
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
