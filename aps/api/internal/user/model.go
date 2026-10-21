package user

import "time"

type User struct {
	ID             int32     `json:"id"`
	TelegramUserID int64     `json:"telegram_user_id"`
	Username       *string   `json:"username"`
	Name           *string   `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateUserParams struct {
	TelegramUserID int64   `json:"telegram_user_id" validate:"required,min=1" minimum:"1" example:"123456789"`
	Username       *string `json:"username" validate:"omitnil,max=64" maxLength:"64" example:"ivan"`
	Name           *string `json:"name" validate:"omitnil,max=128" maxLength:"128" example:"Ivan"`
}

type TelegramIDParams struct {
	TelegramUserID int64 `json:"telegram_user_id" validate:"min=1"`
}
