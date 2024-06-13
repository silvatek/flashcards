package main

import "testing"

func TestQueryParam(t *testing.T) {
	assertQueryParam(t, "/resource?key=value", "key", "value")
	assertQueryParam(t, "/resource", "key", "")
	assertQueryParam(t, "/resource?name=nothing", "key", "")
	assertQueryParam(t, "/resource?a=1&b=2", "a", "1")
	assertQueryParam(t, "/resource?a=1&b=2", "b", "2")
}

func assertQueryParam(t *testing.T, uri string, key string, expectedValue string) {
	value := queryParam(uri, key)
	if value != expectedValue {
		t.Errorf("Unexpected query parameter value `%s` (%s %s %s)", value, uri, key, value)
	}
}

func TestLastPathElement(t *testing.T) {
	assertPathElement(t, "/abc/xyz", "xyz")
	assertPathElement(t, "/123", "123")
	assertPathElement(t, "noslash", "noslash")
	assertPathElement(t, "/abc/xyz?test=1", "xyz")
}

func assertPathElement(t *testing.T, path string, expected string) {
	result := lastPathElement(path)
	if result != expected {
		t.Errorf("lastPathElement returned [%s], expected [%s]", result, expected)
	}
}
