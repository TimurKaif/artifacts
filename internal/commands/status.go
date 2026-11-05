package commands

import "fmt"

func (e *Executor) Status() {
	info, err := e.client.GetCharacter()
	if err != nil {
		fmt.Printf("Не удалось получить статус: %v\n", err)
		return
	}

	fmt.Printf("Персонаж: %s\n", info.Data.Name)
	fmt.Printf("Позиция: X=%d, Y=%d\n", info.Data.X, info.Data.Y)
	fmt.Printf("Жизни: %d\n", info.Data.HP)
}
