package user

import "time"

type User struct {
	ID             int       `json:"id"`
	TelegramUserID int       `json:"telegram_user_id"`
	Username       *string   `json:"username"`
	Name           *string   `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
}
