package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

var ErrEmptyResponse = errors.New("model returned no text")

type Client struct {
	client  openai.Client
	model   string
	timeout time.Duration
}

func NewClient(apiKey, baseURL, model string, timeout time.Duration) *Client {
	return &Client{
		client: openai.NewClient(
			option.WithAPIKey(apiKey),
			option.WithBaseURL(baseURL),
		),
		model:   model,
		timeout: timeout,
	}
}

func (c *Client) Generate(ctx context.Context, instructions, input string) (string, error) {
	requestCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	response, err := c.client.Responses.New(requestCtx, responses.ResponseNewParams{
		Instructions: openai.String(instructions),
		Input:        responses.ResponseNewParamsInputUnion{OfString: openai.String(input)},
		Model:        c.model,
	})
	if err != nil {
		if requestCtx.Err() != nil {
			return "", fmt.Errorf("LLM request: %w", requestCtx.Err())
		}
		return "", fmt.Errorf("LLM request: %w", err)
	}
	if response.Status != responses.ResponseStatusCompleted {
		return "", fmt.Errorf("LLM response: status=%s code=%s reason=%s message=%s",
			response.Status, response.Error.Code, response.IncompleteDetails.Reason, response.Error.Message)
	}

	text := response.OutputText()
	if strings.TrimSpace(text) == "" {
		return "", ErrEmptyResponse
	}
	return text, nil
}
