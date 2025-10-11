package handler

import (
	"net/http"

	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/web"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		web.RenderWithLayout(w, "layout.html", "templates/home.html", map[string]any{
			"Session": sessionData,
		})
	}
}
