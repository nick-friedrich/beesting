package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nick-friedrich/beesting/app/example-api/db"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/views"
	postviews "github.com/nick-friedrich/beesting/app/example-api/views/posts"
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
		views.Layout(postviews.Index(posts), sessionData, "Posts").Render(r.Context(), w)
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
		views.Layout(postviews.Show(post), sessionData, "Post").Render(r.Context(), w)
	}
}

func CreatePostShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)
		if sessionData.UserRole != "admin" {
			Error("You're not authorized to create posts")(w, r)
			return
		}

		views.Layout(postviews.New(), sessionData, "New Post").Render(r.Context(), w)
	}
}

func CreatePostSubmit(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)
		if sessionData.UserRole != "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var input struct {
			Title     string `json:"title"`
			Content   string `json:"content"`
			Author    string `json:"author"`
			Published bool   `json:"published"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		post, err := q.CreatePost(r.Context(), db.CreatePostParams{
			Title:     input.Title,
			Content:   input.Content,
			Author:    input.Author,
			Published: input.Published,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/posts/%s", post.Slug), http.StatusSeeOther)
	}
}
