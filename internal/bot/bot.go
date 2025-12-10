// internal/bot/bot.go
package bot

import (
	"fmt"
	"log"
	"myblog/internal/models"
	"myblog/internal/repo"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	botapi  *tgbotapi.BotAPI
	repo    *repo.Repository
	state   map[int64]string // chat_id → состояние ("waiting_title", "waiting_content")
	pending map[int64]*models.Post
}

func NewBot(token string, repo *repo.Repository) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		botapi:  bot,
		repo:    repo,
		state:   make(map[int64]string),
		pending: make(map[int64]*models.Post),
	}

	log.Printf("Telegram‑бот запущен: @%s", bot.Self.UserName)
	return b, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botapi.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		// Проверяем, админ ли это (по chat_id из .env)
		if strconv.FormatInt(chatID, 10) != os.Getenv("TG_CHAT_ID") {
			b.send(chatID, "Доступ запрещён")
			continue
		}

		switch b.state[chatID] {
		case "waiting_title":
			b.pending[chatID] = &models.Post{Title: update.Message.Text}
			b.state[chatID] = "waiting_content"
			b.send(chatID, "Теперь отправь текст поста (Markdown):")
		case "waiting_content":
			post := b.pending[chatID]
			post.Content = update.Message.Text
			if err := b.repo.CreatePost(post); err != nil {
				b.send(chatID, "Ошибка сохранения :(")
			} else {
				url := fmt.Sprintf("http://localhost:3000/post/%d", post.ID)
				b.send(chatID, fmt.Sprintf("Пост опубликован!\n%s", url))
			}
			delete(b.state, chatID)
			delete(b.pending, chatID)
		default:
			if update.Message.Text == "/new" {
				b.state[chatID] = "waiting_title"
				b.send(chatID, "Отправь заголовок поста:")
			} else {
				b.send(chatID, "Привет! Используй /new чтобы создать пост")
			}
		}
	}
}

func (b *Bot) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.botapi.Send(msg)
}
