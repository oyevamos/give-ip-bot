package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	expectingPassword map[int64]bool
)

func getPublicIP() string {
	resp, err := http.Get("http://api.ipify.org")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(ip)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке файла .env")
	}
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	expectedPassword := os.Getenv("ACCESS_PASSWORD")

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	expectingPassword = make(map[int64]bool)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Text == "/ip" {
			ip := getPublicIP()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, ip)
			bot.Send(msg)
		} else if update.Message.Text == "/password" {
			expectingPassword[update.Message.Chat.ID] = true
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите ваш пароль:")
			bot.Send(msg)
		} else if expectingPassword[update.Message.Chat.ID] {
			if update.Message.Text == expectedPassword {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы успешно авторизованы.")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный пароль.")
				bot.Send(msg)
			}
			expectingPassword[update.Message.Chat.ID] = false // Сброс ожидания ввода пароля
		}
	}
}
