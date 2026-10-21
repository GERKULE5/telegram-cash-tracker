package agent

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

//go:embed prompt.txt
var masterPrompt string

const maxSteps = 8

var ErrStepLimit = errors.New("agent step limit reached")

type Model interface {
	Generate(ctx context.Context, instructions, input string) (string, error)
}

type Executor interface {
	Execute(ctx context.Context, userID int32, action Action) (json.RawMessage, error)
}

type Loop struct {
	model    Model
	executor Executor
}

func NewLoop(model Model, executor Executor) *Loop {
	return &Loop{model: model, executor: executor}
}

func (loop *Loop) Run(ctx context.Context, userID int32, conversation *Context) (string, error) {
	if userID <= 0 {
		return "", errors.New("internal user ID must be positive")
	}
	if conversation == nil {
		return "", errors.New("conversation context is required")
	}
	if conversation.Actions == nil {
		conversation.Actions = make([]ExecutedAction, 0)
	}
	if conversation.CurrentTime.IsZero() {
		return "", errors.New("current time is required in conversation context")
	}

	for step := 0; step < maxSteps; step++ {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		input, err := json.Marshal(conversation)
		if err != nil {
			return "", fmt.Errorf("encode agent context: %w", err)
		}
		output, err := loop.model.Generate(ctx, masterPrompt, string(input))
		if err != nil {
			return "", fmt.Errorf("get next action: %w", err)
		}

		var action Action
		decoder := json.NewDecoder(strings.NewReader(output))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&action); err != nil {
			return "", fmt.Errorf("decode model action: %w", err)
		}
		arguments := strings.TrimSpace(string(action.Arguments))
		if !strings.HasPrefix(arguments, "{") {
			return "", errors.New("action arguments must be a JSON object")
		}

		switch action.Action {
		case ActionReply:
			var reply ReplyArguments
			if err := json.Unmarshal(action.Arguments, &reply); err != nil {
				return "", fmt.Errorf("decode reply arguments: %w", err)
			}
			if strings.TrimSpace(reply.Text) == "" {
				return "", errors.New("reply text must not be empty")
			}
			return reply.Text, nil
		case ActionCreate, ActionList, ActionUpdate, ActionDelete:
			response, executeErr := loop.executor.Execute(ctx, userID, action)
			if err := ctx.Err(); err != nil {
				return "", err
			}
			if executeErr != nil {
				response, err = json.Marshal(ErrorResponse{Error: executeErr.Error()})
				if err != nil {
					return "", fmt.Errorf("encode action error: %w", err)
				}
			} else if !json.Valid(response) {
				return "", fmt.Errorf("action %s returned invalid JSON", action.Action)
			}
			conversation.Actions = append(conversation.Actions, ExecutedAction{
				Action:    action.Action,
				Arguments: action.Arguments,
				Response:  response,
			})
		default:
			return "", fmt.Errorf("unknown action: %q", action.Action)
		}
	}
	return "", ErrStepLimit
}
