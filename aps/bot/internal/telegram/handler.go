package telegram

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"
	_ "time/tzdata"

	"bot/internal/agent"
	"bot/internal/api"
	"bot/internal/llm"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	loop       *agent.Loop
	api        *api.Client
	location   *time.Location
	sessions   map[int64]*session
	sessionsMu sync.Mutex
}

type session struct {
	mu           sync.Mutex
	userID       int32
	conversation agent.Context
}

func NewHandler(loop *agent.Loop, client *api.Client, location *time.Location) *Handler {
	return &Handler{
		loop:     loop,
		api:      client,
		location: location,
		sessions: make(map[int64]*session),
	}
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	message := update.Message
	if message == nil || message.From == nil || strings.TrimSpace(message.Text) == "" {
		return
	}
	if message.Chat.Type != models.ChatTypePrivate {
		return
	}

	h.sessionsMu.Lock()
	current, exists := h.sessions[message.From.ID]
	if !exists {
		current = &session{}
		h.sessions[message.From.ID] = current
	}
	h.sessionsMu.Unlock()

	current.mu.Lock()
	defer current.mu.Unlock()
	if ctx.Err() != nil {
		return
	}

	if current.userID == 0 {
		params := api.UserParams{TelegramUserID: message.From.ID}
		if message.From.Username != "" {
			params.Username = &message.From.Username
		}
		name := []rune(strings.TrimSpace(message.From.FirstName + " " + message.From.LastName))
		if len(name) > 128 {
			name = name[:128]
		}
		if len(name) > 0 {
			nameText := string(name)
			params.Name = &nameText
		}
		user, err := h.api.GetOrCreateUser(ctx, params)
		if err != nil {
			log.Printf("resolve Telegram user: %v", err)
			if _, sendErr := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: message.Chat.ID, Text: "Could not connect to the expense service. Please try again later.",
			}); sendErr != nil {
				log.Printf("send Telegram response: %v", sendErr)
			}
			return
		}
		current.userID = user.ID
	}

	conversation := &current.conversation
	conversation.UserMessage = message.Text
	conversation.CurrentTime = time.Now().In(h.location)
	conversation.Actions = nil
	text, err := h.loop.Run(ctx, current.userID, conversation)

	if err != nil {
		log.Printf("agent processing failed: %v", err)
		text = "Could not finish processing your message. Please try again later."
		if errors.Is(err, context.DeadlineExceeded) {
			text = "The request timed out."
		} else if errors.Is(err, llm.ErrEmptyResponse) {
			text = "The model returned an empty response."
		} else if errors.Is(err, agent.ErrStepLimit) {
			text = "The request exceeded the maximum number of steps. Please clarify what you want to do."
		}
	}
	conversation.History = append(conversation.History, agent.Exchange{
		UserMessage: message.Text,
		Reply:       text,
		Actions:     conversation.Actions,
	})
	if len(conversation.History) > 10 {
		conversation.History = append([]agent.Exchange(nil), conversation.History[len(conversation.History)-10:]...)
	}

	runes := []rune(text)
	for len(runes) > 0 {
		size := min(len(runes), 2000)
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: message.Chat.ID, Text: string(runes[:size]),
		}); err != nil {
			log.Printf("send Telegram response: %v", err)
			return
		}
		runes = runes[size:]
	}
}
