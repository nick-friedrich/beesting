package handler

import (
	"net/http"

	"github.com/nick-friedrich/beesting/app/example-api/pkg/web"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		web.Render(w, "home.html", map[string]any{
			"Title": "Home",
		})
	}
}
