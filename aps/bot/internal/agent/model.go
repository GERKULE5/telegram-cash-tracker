package agent

import (
	"encoding/json"
	"time"
)

type ActionName string

const (
	ActionCreate ActionName = "create"
	ActionList   ActionName = "list"
	ActionUpdate ActionName = "update"
	ActionDelete ActionName = "delete"
	ActionReply  ActionName = "reply"
)

type Action struct {
	Action    ActionName      `json:"action"`
	Arguments json.RawMessage `json:"arguments"`
}

type ExecutedAction struct {
	Action    ActionName      `json:"action"`
	Arguments json.RawMessage `json:"arguments"`
	Response  json.RawMessage `json:"response"`
}

type Context struct {
	UserMessage string           `json:"user_message"`
	CurrentTime time.Time        `json:"current_time"`
	History     []Exchange       `json:"history,omitempty"`
	Actions     []ExecutedAction `json:"actions"`
}

type Exchange struct {
	UserMessage string           `json:"user_message"`
	Reply       string           `json:"reply"`
	Actions     []ExecutedAction `json:"actions,omitempty"`
}

type CreateArguments struct {
	Amount   int32     `json:"amount"`
	Category string    `json:"category"`
	SpentAt  time.Time `json:"spent_at"`
	Comment  *string   `json:"comment,omitempty"`
}

type ListArguments struct {
	Category string     `json:"category,omitempty"`
	From     *time.Time `json:"from,omitempty"`
	To       *time.Time `json:"to,omitempty"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}

type UpdateArguments struct {
	SpendingID int32      `json:"spending_id"`
	Amount     *int32     `json:"amount,omitempty"`
	Category   *string    `json:"category,omitempty"`
	SpentAt    *time.Time `json:"spent_at,omitempty"`
	Comment    *string    `json:"comment,omitempty"`
}

type DeleteArguments struct {
	SpendingID int32 `json:"spending_id"`
}

type ReplyArguments struct {
	Text string `json:"text"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
