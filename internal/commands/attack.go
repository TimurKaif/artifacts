package commands

import "fmt"

func (e *Executor) Attack(participants []string) {
	if len(participants) > 0 {
		fmt.Printf("Атакую монстра вместе с союзниками: %v...\n", participants)
	} else {
		fmt.Println("Атакую монстра на текущей карте...")
	}

	resp, err := e.client.Attack(participants)
	if err != nil {
		e.handleError(err)
		return
	}

	if err := checkResponse(resp); err != nil {
		e.handleError(err)
		return
	}

	fmt.Println("Бой завершен!")
	fmt.Printf("Ответ API: %s\n", resp.Body)
}
