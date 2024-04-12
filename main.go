package main

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oyevamos/give-ip-bot.git/config"
	"github.com/oyevamos/give-ip-bot.git/repository"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
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

	for update := range updates {
		if update.Message == nil {
			continue
		}

		ctx := context.Background()

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		switch update.Message.Command() {
		case "ip":
			session, err := repo.GetSession(ctx, update.Message.Chat.ID)
			if err != nil {
				log.Print(err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Need to authorize, use command /password"))
				continue
			}
			if session.ExpiredAt.Before(time.Now()) {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Need to authorize, use command /password"))
				continue
			}
			ip := getPublicIP()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, ip)
			bot.Send(msg)
		case "password":
			pass := strings.Split(update.Message.Text, "/password ")
			if len(pass) != 2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неправильный формат заполнения пароля")
				bot.Send(msg)
				continue
			}
			if cfg.AccessPassword == pass[1] {
				err = repo.SetSession(ctx, update.Message.Chat.ID)
				if err != nil {
					log.Print(err)

				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный пароль")
				bot.Send(msg)
			}
		}
	}
}
