package handler

import (
	"net/http"

	"github.com/nick-friedrich/beesting/app/example-api/pkg/web"
)

func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		web.RenderWithLayout(w, "layout.html", "templates/404.html", map[string]any{})
	}
}
