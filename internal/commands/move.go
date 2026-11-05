package commands

import "fmt"

func (e *Executor) Move(x, y int) {
	fmt.Printf("Перемещаюсь в точку (%d, %d)...\n", x, y)

	result, err := e.client.Move(x, y)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Успешно переместился! Ура! Ура!\n")
	fmt.Printf("Ответ API: %s\n", result)
}
