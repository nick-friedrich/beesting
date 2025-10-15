package handler

import (
	"net/http"

	"github.com/nick-friedrich/beesting/app/example-app/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-app/views"
)

func Error(error string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		w.WriteHeader(http.StatusInternalServerError)
		views.Layout(
			views.ErrorView(error),
			sessionData,
			"Oops. Something went wrong",
		).Render(r.Context(), w)

	}
}
