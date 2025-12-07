package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"myblog/internal/handlers"
	"myblog/internal/repo"
	"myblog/internal/templates"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PageData struct {
	SiteTitle string
	IsAdmin   bool
	// сюда потом можно добавить Description, Author и т.д.
}

func main() {
	godotenv.Load()

	// БД
	db, err := sql.Open("postgres", os.ExpandEnv("host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD sslmode=$DB_SSLMODE"))
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Репозиторий
	repo := repo.NewRepository(db)

	// Шаблоны
	tpl, err := template.ParseFS(templates.FS, "index.html", "post.html", "admin_login.html")
	if err != nil {
		log.Fatal("Шаблоны не найдены:", err)
	}

	// Хендлеры
	h := handlers.NewHandler(repo, tpl)

	// Маршруты
	// === АДМИНКА ===
	http.HandleFunc("/admin", h.AdminLoginPage)
	http.HandleFunc("/login", h.LoginHandler)
	http.HandleFunc("/logout", h.LogoutHandler)

	// Защищённые маршруты (только админ)
	http.HandleFunc("/new", h.AdminOnly(h.NewPostHandler))
	http.HandleFunc("/create", h.AdminOnly(h.CreatePostHandler))
	http.HandleFunc("/edit/", h.AdminOnly(h.EditPostHandler))
	http.HandleFunc("/update/", h.AdminOnly(h.UpdatePostHandler))
	http.HandleFunc("/delete/", h.AdminOnly(h.DeletePostHandler))

	// Публичные
	http.HandleFunc("/post/", h.PostHandler)
	http.HandleFunc("/", h.HomeHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Printf("Блог запущен → http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
