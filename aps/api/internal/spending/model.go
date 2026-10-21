package spending

import "time"

type Spending struct {
	ID        int32      `json:"id"`
	UserID    int32      `json:"user_id"`
	Amount    int32      `json:"amount"`
	Category  Category   `json:"category"`
	SpentAt   time.Time  `json:"spent_at"`
	Comment   *string    `json:"comment"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type SpendingWithUsername struct {
	Spending
	Username *string `json:"username"`
}

type CreateParams struct {
	UserID   int32     `json:"user_id" validate:"required,min=1" minimum:"1" example:"1"`
	Amount   int32     `json:"amount" validate:"required,min=1" minimum:"1" example:"350"`
	Category Category  `json:"category" validate:"required,oneof=food transport housing entertainment health purchases other" enums:"food,transport,housing,entertainment,health,purchases,other"`
	SpentAt  time.Time `json:"spent_at" validate:"required" example:"2026-10-17T12:00:00+07:00"`
	Comment  *string   `json:"comment" example:"Coffee"`
}

type UpdateParams struct {
	Amount   *int32     `json:"amount" validate:"omitnil,min=1" minimum:"1" example:"450"`
	Category *Category  `json:"category" validate:"omitnil,oneof=food transport housing entertainment health purchases other" enums:"food,transport,housing,entertainment,health,purchases,other"`
	SpentAt  *time.Time `json:"spent_at" validate:"omitnil,notzerotime" example:"2026-10-17T12:00:00+07:00"`
	Comment  *string    `json:"comment" example:"Coffee"`
}

type ListParams struct {
	UserID   int32    `json:"user_id" validate:"min=0"`
	Category Category `json:"category" validate:"omitempty,oneof=food transport housing entertainment health purchases other"`
	From     *time.Time
	To       *time.Time
	Limit    int `json:"limit" validate:"min=1,max=100"`
	Offset   int `json:"offset" validate:"min=0"`
}

type IDParams struct {
	ID int32 `json:"id" validate:"min=1"`
}

type UserIDParams struct {
	UserID int32 `json:"user_id" validate:"min=1"`
}

type PeriodParams struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to" validate:"gtfield=From"`
}
