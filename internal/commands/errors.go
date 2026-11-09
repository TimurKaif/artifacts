package commands

import "fmt"

type CooldownError struct {
	Seconds int
}

func (e *CooldownError) Error() string {
	return fmt.Sprintf("персонаж в кулдауне, осталось %d сек", e.Seconds)
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("ошибка API (%d): %s", e.StatusCode, e.Body)
}
