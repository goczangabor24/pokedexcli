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
		fmt.Printf("Your command was: %v\n", cleanInput(scanner.Text())[0])
	}
}
