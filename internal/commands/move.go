package commands

import "fmt"

func (e *Executor) Move(x, y int) {
	fmt.Printf("Я пошёл в точку (%d, %d)...\n", x, y)

	resp, err := e.client.Move(x, y)
	if err != nil {
		e.handleError(err)
		return
	}

	if err := checkResponse(resp); err != nil {
		e.handleError(err)
		return
	}

	fmt.Println("Успешно переместился! Ура! Ура!")
	fmt.Printf("Ответ API: %s\n", resp.Body)
}
