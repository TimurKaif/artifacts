package commands

import "artifacts/internal/api"

type Executor struct {
	client *api.Client
}

func NewExecutor(client *api.Client) *Executor {
	return &Executor{client: client}
}
