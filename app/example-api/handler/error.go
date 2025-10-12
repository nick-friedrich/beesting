package handler

import (
	"net/http"

	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/web"
)

func Error(error string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		w.WriteHeader(http.StatusInternalServerError)
		web.RenderWithLayoutAndSession(w, "layout.html", "templates/error.html", map[string]any{
			"error": error,
		}, sessionData)
	}
}
