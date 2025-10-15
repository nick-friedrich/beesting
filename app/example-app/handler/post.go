package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nick-friedrich/beesting/app/example-app/db"
	"github.com/nick-friedrich/beesting/app/example-app/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-app/pkg/slug"
	"github.com/nick-friedrich/beesting/app/example-app/views"
	postviews "github.com/nick-friedrich/beesting/app/example-app/views/posts"
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
		views.Layout(postviews.Index(posts, sessionData, r), sessionData, "Posts").Render(r.Context(), w)
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
		views.Layout(postviews.Show(post, sessionData, r), sessionData, "Post").Render(r.Context(), w)
	}
}

func CreatePostShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)
		if sessionData.UserRole != "admin" {
			Error("You're not authorized to create posts")(w, r)
			return
		}

		views.Layout(postviews.New(r), sessionData, "New Post").Render(r.Context(), w)
	}
}

func CreatePostSubmit(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)
		if sessionData.UserRole != "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		author := r.FormValue("author")
		slugValue := r.FormValue("slug")
		published := r.FormValue("published") == "on" // Checkbox sends "on" when checked

		if title == "" || content == "" || author == "" || slugValue == "" {
			http.Error(w, "Title, content, author, and slug are required", http.StatusBadRequest)
			return
		}

		// Validate slug
		if err := slug.Validate(slugValue); err != nil {
			http.Error(w, fmt.Sprintf("Invalid slug: %s", err.Error()), http.StatusBadRequest)
			return
		}

		// Check if slug already exists
		existingPost, err := q.GetPostBySlug(r.Context(), slugValue)
		if err == nil && existingPost.ID != 0 {
			http.Error(w, "A post with this slug already exists", http.StatusBadRequest)
			return
		}

		post, err := q.CreatePost(r.Context(), db.CreatePostParams{
			Title:     title,
			Slug:      slugValue,
			Content:   content,
			Author:    author,
			Published: published,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/posts/%s", post.Slug), http.StatusSeeOther)
	}
}

func EditPostShow(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)
		if sessionData.UserRole != "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		post, err := q.GetPost(r.Context(), postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		views.Layout(postviews.Edit(post, r), sessionData, "Edit Post").Render(r.Context(), w)
	}
}

func EditPostSubmit(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)
		if sessionData.UserRole != "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		author := r.FormValue("author")
		slugValue := r.FormValue("slug")
		published := r.FormValue("published") == "on"

		if title == "" || content == "" || author == "" || slugValue == "" {
			http.Error(w, "Title, content, author, and slug are required", http.StatusBadRequest)
			return
		}

		// Validate slug
		if err := slug.Validate(slugValue); err != nil {
			http.Error(w, fmt.Sprintf("Invalid slug: %s", err.Error()), http.StatusBadRequest)
			return
		}

		// Check if slug already exists (excluding current post)
		existingPost, err := q.GetPostBySlug(r.Context(), slugValue)
		if err == nil && existingPost.ID != 0 && existingPost.ID != postID {
			http.Error(w, "A post with this slug already exists", http.StatusBadRequest)
			return
		}

		post, err := q.UpdatePost(r.Context(), db.UpdatePostParams{
			ID:        postID,
			Title:     title,
			Slug:      slugValue,
			Content:   content,
			Author:    author,
			Published: published,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/posts/%s", post.Slug), http.StatusSeeOther)
	}
}

func DeletePostWeb(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)
		if sessionData.UserRole != "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		err = q.DeletePost(r.Context(), postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}
