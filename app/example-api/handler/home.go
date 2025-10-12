package handler

import (
	"net/http"

	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/views"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		views.Layout(
			views.Home(sessionData),
			sessionData,
			"Welcome to the home page",
		).Render(r.Context(), w)

		// web.RenderWithLayout(w, "layout.html", "templates/home.html", map[string]any{
		// 	"Session": sessionData,
		// })
	}
}
