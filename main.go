package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())[0]

		c, ok := commands[command]
		if !ok {
			fmt.Println("Unknown value")
			continue
		}
		c.callback()
	}
}
