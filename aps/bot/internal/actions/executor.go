package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"bot/internal/agent"
	"bot/internal/api"
)

type Executor struct {
	client *api.Client
}

func NewExecutor(client *api.Client) *Executor {
	return &Executor{client: client}
}

func (e *Executor) Execute(ctx context.Context, userID int32, action agent.Action) (json.RawMessage, error) {
	if userID <= 0 {
		return nil, errors.New("internal user ID must be positive")
	}
	switch action.Action {
	case agent.ActionCreate:
		var args agent.CreateArguments
		if err := json.Unmarshal(action.Arguments, &args); err != nil {
			return nil, err
		}
		return e.client.CreateSpending(ctx, userID, args)
	case agent.ActionList:
		args := agent.ListArguments{Limit: 20}
		if err := json.Unmarshal(action.Arguments, &args); err != nil {
			return nil, err
		}
		return e.client.ListSpendings(ctx, userID, args)
	case agent.ActionUpdate:
		var args agent.UpdateArguments
		if err := json.Unmarshal(action.Arguments, &args); err != nil {
			return nil, err
		}
		if args.SpendingID <= 0 {
			return nil, errors.New("spending_id must be positive")
		}
		return e.client.UpdateSpending(ctx, userID, args)
	case agent.ActionDelete:
		var args agent.DeleteArguments
		if err := json.Unmarshal(action.Arguments, &args); err != nil {
			return nil, err
		}
		if args.SpendingID <= 0 {
			return nil, errors.New("spending_id must be positive")
		}
		return e.client.DeleteSpending(ctx, userID, args.SpendingID)
	default:
		return nil, fmt.Errorf("unsupported action: %s", action.Action)
	}
}
