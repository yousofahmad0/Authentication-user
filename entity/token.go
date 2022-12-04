package entity

import (
	"time"
)

type Token struct {
	UserId      uint      `json:"user_id"`
	Token       string    `json:"token"`
	ActivatedAt time.Time `json:"activated_at"`
	ExpiredAt   time.Time `json:"expired_at"`
	Type        string    `json:"type"`
}
