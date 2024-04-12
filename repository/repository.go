package repository

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/oyevamos/give-ip-bot.git/config"
	"github.com/oyevamos/give-ip-bot.git/domain"
	"time"
	"xorm.io/xorm"
)

type WeatherRepository struct {
	Engine *xorm.Engine
}

func NewWeather(cfg config.Postgres) (*WeatherRepository, error) {
	db, err := connectDB(cfg)
	if err != nil {
		return nil, err
	}
	return &WeatherRepository{
		Engine: db,
	}, nil
}
func connectDB(cfg config.Postgres) (*xorm.Engine, error) {
	connectURL := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.DbUser, cfg.DbPassword, cfg.DbName)

	db, err := xorm.NewEngine("postgres", connectURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (w *WeatherRepository) SetSession(ctx context.Context, chatId int64) error {
	expiredAt := time.Now().Add(72 * time.Hour)
	has, err := w.Engine.Get(&domain.Session{
		ChatId: chatId,
	})
	if err != nil {
		return err
	}
	session := &domain.Session{
		ChatId:    chatId,
		ExpiredAt: expiredAt,
	}
	if has {
		_, err = w.Engine.Update(session)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = w.Engine.Insert(session)
	return err
}

func (w *WeatherRepository) GetSession(ctx context.Context, chatId int64) (*domain.Session, error) {
	//_, err := w.Engine.Delete(&domain.Weather{Date: date})
	session := &domain.Session{
		ChatId: chatId,
	}
	has, err := w.Engine.Get(session)
	if !has {
		return nil, errors.New("session not found")
	}

	return session, err
}
