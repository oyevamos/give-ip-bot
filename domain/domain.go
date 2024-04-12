package domain

import "time"

type Session struct {
	ChatId    int64
	ExpiredAt time.Time
}
