package main

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oyevamos/give-ip-bot.git/config"
	"github.com/oyevamos/give-ip-bot.git/repository"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	cfg := config.LoadConfig()

	repo, err := repository.NewWeather(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatal(err)
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
		ctx := context.Background()

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Text == "/ip" {
			session, err := repo.GetSession(ctx, update.Message.Chat.ID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Need to authorize, use command /password"))
				continue
			}
			if session.ExpiresAt.After(time.Now()) {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Need to authorize, use command /password"))
				continue
			}
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
