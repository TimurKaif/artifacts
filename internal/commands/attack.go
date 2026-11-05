package commands

import "fmt"

func (e *Executor) Attack(participants []string) {
	if len(participants) > 0 {
		fmt.Printf("Атакую монстра вместе с союзниками: %v...\n", participants)
	} else {
		fmt.Println("Атакую монстра на текущей карте...")
	}

	result, err := e.client.Attack(participants)
	if err != nil {
		fmt.Printf("Ошибка атаки: %v\n", err)
		return
	}

	fmt.Printf("Бой завершен!\n")
	fmt.Printf("Ответ API: %s\n", result)
}
