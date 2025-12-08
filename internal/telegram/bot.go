// internal/telegram/bot.go
package telegram

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

var (
	botToken    = ""
	chatID      = ""
	initialized = false
)

func init() {
	_ = godotenv.Load()

	botToken = os.Getenv("TG_BOT_TOKEN")
	chatID = os.Getenv("TG_CHAT_ID")

	if botToken == "" || chatID == "" {
		log.Println("Telegram: TG_BOT_TOKEN или TG_CHAT_ID не заданы — уведомления отключены")
	} else {
		initialized = true
		log.Println("Telegram бот инициализирован")
	}
}

func Send(message string) {
	if !initialized {
		return
	}

	text := url.QueryEscape(message)
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=HTML", botToken, chatID, text)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Telegram: ошибка отправки: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Println("Telegram: сообщение отправлено")
	} else {
		log.Printf("Telegram: ошибка %d", resp.StatusCode)
	}
}
