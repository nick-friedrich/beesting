package handler

import (
	"net/http"

	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/views"
)

func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		w.WriteHeader(http.StatusNotFound)
		views.Layout(
			views.NotFound(),
			sessionData,
			"Not found",
		).Render(r.Context(), w)

	}
}
