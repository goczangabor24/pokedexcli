package main

import (
	"time"

	"github.com/goczangabor24/pokedexcli/internal/pokecache"
)

func main() {
	cfg := &config{
		cache: pokecache.NewCache(5 * time.Minute),

		commands: map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
			"map": {
				name:        "map",
				description: "Displays the names of 20 location areas in the Pokemon world, goes to the next 20 on each consecutive call",
				callback:    commandMap,
			},
			"mapb": {
				name:        "mapb",
				description: "Goes back to the previous 20 location areas in the Pokemon world",
				callback:    commandMapb,
			},
		},
	}

	startRepl(cfg)
}
