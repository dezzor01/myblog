package handlers

import (
	"net/http"
	"strconv"

	"myblog/internal/services"
)

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/post/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	post, err := h.repo.GetPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]any{
		"Title":     post.Title,
		"HTML":      services.RenderMarkdown(post.Content),
		"ID":        post.ID,
		"IsAdmin":   h.isAdmin(r),
		"CreatedAt": post.CreatedAt,
	}

	if err := h.tpl.ExecuteTemplate(w, "post.html", data); err != nil {
		http.Error(w, "Ошибка рендера", http.StatusInternalServerError)
	}
}
