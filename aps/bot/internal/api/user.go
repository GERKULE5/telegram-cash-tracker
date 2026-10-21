package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type User struct {
	ID             int32 `json:"id"`
	TelegramUserID int64 `json:"telegram_user_id"`
}

type UserParams struct {
	TelegramUserID int64   `json:"telegram_user_id"`
	Username       *string `json:"username,omitempty"`
	Name           *string `json:"name,omitempty"`
}

func (c *Client) GetUser(ctx context.Context, telegramID int64) (User, error) {
	query := url.Values{"telegram_user_id": {strconv.FormatInt(telegramID, 10)}}
	data, err := c.request(ctx, http.MethodGet, "/users", query, nil)
	if err != nil {
		return User{}, err
	}
	var result User
	if err := json.Unmarshal(data, &result); err != nil {
		return User{}, fmt.Errorf("decode user: %w", err)
	}
	if result.ID <= 0 || result.TelegramUserID != telegramID {
		return User{}, errors.New("API returned an invalid user")
	}
	return result, nil
}

func (c *Client) GetOrCreateUser(ctx context.Context, params UserParams) (User, error) {
	result, err := c.GetUser(ctx, params.TelegramUserID)
	if err == nil {
		return result, nil
	}
	var apiErr *HTTPError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusNotFound {
		return User{}, err
	}
	data, err := c.request(ctx, http.MethodPost, "/users", nil, params)
	if err != nil {
		if errors.As(err, &apiErr) && apiErr.Status == http.StatusConflict {
			return c.GetUser(ctx, params.TelegramUserID)
		}
		return User{}, err
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return User{}, fmt.Errorf("decode created user: %w", err)
	}
	if result.ID <= 0 || result.TelegramUserID != params.TelegramUserID {
		return User{}, errors.New("API returned an invalid created user")
	}
	return result, nil
}
