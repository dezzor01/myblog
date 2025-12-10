// internal/handlers/handlers.go
package handlers

import (
	"html/template"
	"log"
	"myblog/internal/config"
	"myblog/internal/models"
	"myblog/internal/repo"
	"myblog/internal/services"
	"net/http"
	"strconv"
)

type Handler struct {
	repo *repo.Repository
	tpl  *template.Template
	cfg  *config.Config
}

func NewHandler(repo *repo.Repository, tpl *template.Template, cfg *config.Config) *Handler {
	return &Handler{
		repo: repo,
		tpl:  tpl,
		cfg:  cfg,
	}
}

func (h *Handler) isAdmin(r *http.Request) bool {
	cookie, _ := r.Cookie("admin_session")
	return cookie != nil && cookie.Value == h.cfg.AdminPassword
}

// Главная
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	posts, err := h.repo.GetAllPosts()
	if err != nil {
		log.Printf("Детали ошибки БД: %v", err) // ← добавь это
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"IsAdmin":   h.isAdmin(r),
		"SiteTitle": h.cfg.SiteTitle,
		"Posts":     posts,
	}

	if err := h.tpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
	}
}

// ← ДОБАВЛЕННЫЙ МЕТОД — БЕЗ НЕГО НЕ РАБОТАЕТ /post/1 !
func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/post/"):]
	id, _ := strconv.Atoi(idStr)
	if id < 1 {
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
		"SiteTitle": h.cfg.SiteTitle,
		"CreatedAt": post.CreatedAt,
	}

	if err := h.tpl.ExecuteTemplate(w, "post.html", data); err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
	}
}

// Новая запись
func (h *Handler) NewPostHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"NewPost":   true,
		"IsAdmin":   true,
		"SiteTitle": h.cfg.SiteTitle,
	}
	h.tpl.ExecuteTemplate(w, "index.html", data)
}

// Создание поста
func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
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

// Редактирование
func (h *Handler) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/edit/"):]
	id, _ := strconv.Atoi(idStr)
	post, err := h.repo.GetPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]any{
		"EditPost":  true,
		"Post":      post,
		"IsAdmin":   true,
		"SiteTitle": h.cfg.SiteTitle,
	}
	h.tpl.ExecuteTemplate(w, "index.html", data)
}

// Обновление
func (h *Handler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
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

// Удаление
func (h *Handler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
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
