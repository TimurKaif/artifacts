package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/TimurKaif/artifacts/internal/transport"
)

const statusCodeCooldown = 499

const defaultCooldownSeconds = 5

type Executor struct {
	client *transport.Client
}

func NewExecutor(client *transport.Client) *Executor {
	return &Executor{client: client}
}

func checkResponse(resp *transport.Response) error {
	if resp.StatusCode == 200 {
		return nil
	}

	if resp.StatusCode == statusCodeCooldown {
		return cooldownFrom(resp)
	}

	return &APIError{StatusCode: resp.StatusCode, Body: string(resp.Body)}
}

func cooldownFrom(resp *transport.Response) error {
	seconds := 0

	var body struct {
		Cooldown Cooldown `json:"cooldown"`
		Data     struct {
			Cooldown Cooldown `json:"cooldown"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &body); err == nil {
		switch {
		case body.Cooldown.RemainingSeconds > 0:
			seconds = body.Cooldown.RemainingSeconds
		case body.Data.Cooldown.RemainingSeconds > 0:
			seconds = body.Data.Cooldown.RemainingSeconds
		}
	}

	if seconds <= 0 {
		seconds = defaultCooldownSeconds
	}
	return &CooldownError{Seconds: seconds}
}

func (e *Executor) handleError(err error) {
	var cooldownErr *CooldownError
	if errors.As(err, &cooldownErr) {
		fmt.Printf("%s. Кулдаун...\n", cooldownErr.Error())
		return
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("Ошибка: %s\n", apiErr.Error())
		return
	}

	fmt.Printf("Ошибка: %v\n", err)
}
