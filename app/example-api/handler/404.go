package handler

import (
	"net/http"

	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/web"
)

func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionManager := session.NewSessionManager()
		sessionData, _ := sessionManager.GetSession(r)

		w.WriteHeader(http.StatusNotFound)
		web.RenderWithLayoutAndSession(w, "layout.html", "templates/404.html", map[string]any{}, sessionData)
	}
}
