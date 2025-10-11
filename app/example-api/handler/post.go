package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nick-friedrich/beesting/app/example-api/db"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/web"
)

func ShowPosts(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := q.ListPosts(r.Context(), db.ListPostsParams{
			Limit:  10,
			Offset: 0,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionData, _ := session.Default.GetSession(r)
		web.RenderWithLayoutAndSession(w, "layout.html", "templates/posts/index.html", map[string]any{
			"Posts": posts,
		}, sessionData)
	}
}

func ShowPost(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		post, err := q.GetPostBySlug(r.Context(), chi.URLParam(r, "slug"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionData, _ := session.Default.GetSession(r)
		web.RenderWithLayoutAndSession(w, "layout.html", "templates/posts/show.html", map[string]any{
			"Post": post,
		}, sessionData)
	}
}
