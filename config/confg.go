package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Postgres struct {
	Host       string
	Port       int
	DbName     string
	DbPassword string
	DbUser     string
}

type Config struct {
	Postgres         Postgres
	TelegramBotToken string
	AccessPassword   string
}

//POSTGRES_DB=give-ip-bot
//POSTGRES_USER=user
//POSTGRES_PASSWORD=12345
//POSTGRES_HOST=localhost
//POSTGRES_PORT=5432

func LoadConfig() Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке файла .env")
	}
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		panic(err)
	}
	return Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		AccessPassword:   os.Getenv("ACCESS_PASSWORD"),
		Postgres: Postgres{
			Host:       os.Getenv("POSTGRES_HOST"),
			Port:       port,
			DbName:     os.Getenv("POSTGRES_DB"),
			DbUser:     os.Getenv("POSTGRES_USER"),
			DbPassword: os.Getenv("POSTGRES_PASSWORD"),
		},
	}
}
