package handlers

import (
	"net/http"
	"strings"
)

type History struct {
	cookieName string
	entries    []string
}

func getHistory(cookieName string, r *http.Request) History {
	ctx := requestContext(r)
	var history History
	history.cookieName = cookieName
	current, err := r.Cookie(cookieName)
	if err != http.ErrNoCookie {
		history.entries = strings.Split(current.Value, "|")
		logs.Debug(ctx, "Loaded  history: %v", history.entries)
	}
	return history
}

func (h *History) setCookie(w http.ResponseWriter) {
	var value string
	for _, entry := range h.entries {
		if len(value) > 0 {
			value = value + "|"
		}
		value = value + entry
	}
	cookie := http.Cookie{
		Name:     h.cookieName,
		Value:    value,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)
}

func (h *History) push(entry string) {
	// Create a new list with just the new entry
	updated := []string{entry}

	// Then add the existing entries
	for count, val := range h.entries {
		// Don't add duplicates
		if val != entry {
			updated = append(updated, val)
		}
		// List has a maximum length
		if count >= 3 {
			break
		}
	}

	h.entries = updated
}
