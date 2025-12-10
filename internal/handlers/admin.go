package handlers

import (
	"net/http"
	"os"

	"myblog/internal/auth"
)

func (h *Handler) AdminLoginPage(w http.ResponseWriter, r *http.Request) {
	if auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	h.tpl.ExecuteTemplate(w, "admin_login.html", nil)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	if r.FormValue("password") == os.Getenv("ADMIN_PASSWORD") {
		auth.SetAuthCookie(w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Неверный пароль", http.StatusUnauthorized)
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		auth.ClearAuthCookie(w)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
