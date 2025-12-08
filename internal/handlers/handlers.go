package handlers

import (
	"html/template"
	"myblog/internal/config"
	"myblog/internal/models"
	"myblog/internal/repo"
	"net/http"
	"strconv"
)

const MySiteTitle = "Trash"

type Handler struct {
	repo *repo.Repository
	tpl  *template.Template
	cfg  *config.Config
}

func (h *Handler) isAdmin(r *http.Request) bool {
	cookie, _ := r.Cookie("admin_session")
	return cookie != nil && cookie.Value == h.cfg.AdminPassword
}

func NewHandler(repo *repo.Repository, tpl *template.Template, cfg *config.Config) *Handler {
	return &Handler{
		repo: repo,
		tpl:  tpl,
		cfg:  cfg,
	}
}
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	posts, err := h.repo.GetAllPosts()
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"IsAdmin":   h.isAdmin(r),
		"SiteTitle": h.cfg.SiteTitle, // ← ЭТОЙ СТРОКИ НЕ БЫЛО!
		"Posts":     posts,
	}

	if err := h.tpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
	}
}

// Новая запись — форма
func (h *Handler) NewPostHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
}

// Создание поста
func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")
	if title == "" || content == "" {
		http.Error(w, "Заполни поля", http.StatusBadRequest)
		return
	}

	post := &models.Post{Title: title, Content: content}
	if err := h.repo.CreatePost(post); err != nil {
		http.Error(w, "Ошибка сохранения", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(post.ID), http.StatusSeeOther)
}

// Редактирование — форма
func (h *Handler) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	idStr := r.URL.Path[len("/edit/"):]
	id, _ := strconv.Atoi(idStr)
	post, err := h.repo.GetPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	h.tpl.ExecuteTemplate(w, "index.html", map[string]any{
		"EditPost": true,
		"Post":     post,
		"IsAdmin":  true,
	})
}

// Обновление поста
func (h *Handler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	idStr := r.URL.Path[len("/update/"):]
	id, _ := strconv.Atoi(idStr)
	title := r.FormValue("title")
	content := r.FormValue("content")

	post := &models.Post{ID: id, Title: title, Content: content}
	if err := h.repo.UpdatePost(post); err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(id), http.StatusSeeOther)
}

// Удаление поста
func (h *Handler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	idStr := r.URL.Path[len("/delete/"):]
	id, _ := strconv.Atoi(idStr)
	h.repo.DeletePost(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.isAdmin(r) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
