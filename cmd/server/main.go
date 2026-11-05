package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"artifacts/config"
	"artifacts/internal/api"
	"artifacts/internal/commands"
)

func main() {
	cfg := config.NewConfig()
	client := api.NewClient(cfg)
	executor := commands.NewExecutor(client)

	fmt.Println("Введите 'help' для списка команд")
	fmt.Println()

	executor.Status()
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := parts[0]

		switch cmd {
		case "move":
			if len(parts) != 3 {
				fmt.Println("Использование: move X Y")
				fmt.Println("   Пример: move 4 -1")
				continue
			}

			x, err1 := strconv.Atoi(parts[1])
			y, err2 := strconv.Atoi(parts[2])
			if err1 != nil || err2 != nil {
				fmt.Println("X и Y должны быть числами")
				continue
			}

			executor.Move(x, y)

		case "attack":
			executor.Attack(nil)

		case "status", "info":
			executor.Status()

		case "help":
			printHelp()

		case "exit", "quit", "q":
			fmt.Println("До встречи!")
			return

		default:
			fmt.Printf("Неизвестная команда: %s\n", cmd)
			fmt.Println("   Введите 'help' для справки")
		}

		fmt.Println()
	}
}

func printHelp() {
	fmt.Println("Доступные команды:")
	fmt.Println("  move X Y    - Переместиться в точку (X, Y)")
	fmt.Println("  status      - Показать текущий статус персонажа")
	fmt.Println("  attack      - Атаковать")
	fmt.Println("  help        - Показать эту справку")
	fmt.Println("  exit        - Выйти из программы")
}
