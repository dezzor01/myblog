package auth

import (
	"net/http"
	"os"
)

const cookieName = "admin_session"

func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}
	return cookie.Value == os.Getenv("ADMIN_PASSWORD")
}

func SetAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    os.Getenv("ADMIN_PASSWORD"),
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60, // 30 дней
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
