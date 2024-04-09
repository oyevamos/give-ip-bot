package domain

import "time"

type Session struct {
	ChatId    int64
	ExpiresAt time.Time
}
