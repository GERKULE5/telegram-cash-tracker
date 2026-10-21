package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bot/internal/agent"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: timeout}}
}

type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("expenses API: status=%d: %s", e.Status, e.Message)
}

func (c *Client) request(ctx context.Context, method, path string, query url.Values, body any) (json.RawMessage, error) {
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encode API request: %w", err)
		}
	}
	address := c.baseURL + path
	if len(query) > 0 {
		address += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, address, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create API request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("expenses API request (result unconfirmed): %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read API response (result unconfirmed): %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var response struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(data, &response); err != nil || response.Error == "" {
			response.Error = http.StatusText(resp.StatusCode)
		}
		return nil, &HTTPError{Status: resp.StatusCode, Message: response.Error}
	}
	if resp.StatusCode == http.StatusNoContent {
		return json.RawMessage(`{}`), nil
	}
	if !json.Valid(data) {
		return nil, errors.New("expenses API returned invalid JSON (result unconfirmed)")
	}
	return data, nil
}

func (c *Client) CreateSpending(ctx context.Context, userID int32, args agent.CreateArguments) (json.RawMessage, error) {
	body := struct {
		UserID int32 `json:"user_id"`
		agent.CreateArguments
	}{UserID: userID, CreateArguments: args}
	return c.request(ctx, http.MethodPost, "/spendings", nil, body)
}

func (c *Client) ListSpendings(ctx context.Context, userID int32, args agent.ListArguments) (json.RawMessage, error) {
	query := url.Values{
		"user_id": {strconv.FormatInt(int64(userID), 10)},
		"limit":   {strconv.Itoa(args.Limit)},
		"offset":  {strconv.Itoa(args.Offset)},
	}
	if args.Category != "" {
		query.Set("category", args.Category)
	}
	if args.From != nil {
		query.Set("from", args.From.Format(time.RFC3339))
	}
	if args.To != nil {
		query.Set("to", args.To.Format(time.RFC3339))
	}
	items, err := c.request(ctx, http.MethodGet, "/spendings", query, nil)
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Items json.RawMessage `json:"items"`
	}{Items: items})
}

func (c *Client) UpdateSpending(ctx context.Context, userID int32, args agent.UpdateArguments) (json.RawMessage, error) {
	query := url.Values{"user_id": {strconv.FormatInt(int64(userID), 10)}}
	body := struct {
		Amount   *int32     `json:"amount,omitempty"`
		Category *string    `json:"category,omitempty"`
		SpentAt  *time.Time `json:"spent_at,omitempty"`
		Comment  *string    `json:"comment,omitempty"`
	}{Amount: args.Amount, Category: args.Category, SpentAt: args.SpentAt, Comment: args.Comment}
	return c.request(ctx, http.MethodPatch, fmt.Sprintf("/spendings/%d", args.SpendingID), query, body)
}

func (c *Client) DeleteSpending(ctx context.Context, userID, spendingID int32) (json.RawMessage, error) {
	query := url.Values{"user_id": {strconv.FormatInt(int64(userID), 10)}}
	_, err := c.request(ctx, http.MethodDelete, fmt.Sprintf("/spendings/%d", spendingID), query, nil)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(`{"deleted":true}`), nil
}
