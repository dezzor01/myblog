package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"myblog/internal/bot"
	"myblog/internal/config"
	"myblog/internal/handlers"
	"myblog/internal/repo"
	"myblog/internal/telegram"
	"myblog/internal/templates"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type PageData struct {
	SiteTitle string
	IsAdmin   bool
	// сюда потом можно добавить Description, Author и т.д.
}

func main() {
	cfg := config.Load()

	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBName,
		cfg.DBUser, cfg.DBPassword, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
		telegram.Send("Блог не запустился! Ошибка подключения к БД: " + err.Error())
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Не удалось пингануть БД:", err)
	}
	fmt.Println("Подключено к PostgreSQL")

	repo := repo.NewRepository(db)
	tpl, err := template.ParseFS(templates.FS, "index.html", "post.html", "admin_login.html")
	if err != nil {
		log.Fatal("Ошибка загрузки шаблонов:", err)
	}

	h := handlers.NewHandler(repo, tpl, cfg) // ← передаём конфиг

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

	// подключаем тг-бот
	go func() {
		bot, err := bot.NewBot(os.Getenv("TG_BOT_TOKEN"), repo)
		if err != nil {
			log.Fatal("Не удалось запустить Telegram‑бота:", err)
		}
		bot.Start()
	}()
	fmt.Printf("Блог запущен → http://localhost:%s\n", port)
	telegram.Send(fmt.Sprintf("Блог запущен!\nhttp://localhost:%s", port))

	// Защита от паники
	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("ПАНИКА! Блог упал: %v", r)
			log.Println(errMsg)
			telegram.Send(errMsg)
			panic(r)
		}
	}()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
