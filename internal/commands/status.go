package commands

import (
	"encoding/json"
	"fmt"
)

func (e *Executor) Status() {
	resp, err := e.client.GetCharacter()
	if err != nil {
		e.handleError(err)
		return
	}

	if err := checkResponse(resp); err != nil {
		e.handleError(err)
		return
	}

	var wrapper struct {
		Data CharacterInfo `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &wrapper); err != nil {
		e.handleError(fmt.Errorf("не удалось разобрать ответ: %w", err))
		return
	}

	fmt.Printf("Персонаж: %s\n", wrapper.Data.Name)
	fmt.Printf("Позиция: X=%d, Y=%d\n", wrapper.Data.X, wrapper.Data.Y)
	fmt.Printf("Жизни: %d\n", wrapper.Data.HP)
}
