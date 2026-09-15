package main

import (
	"time"

	"github.com/goczangabor24/pokedexcli/internal/pokecache"
)

func main() {
	cfg := &config{
		cache:          pokecache.NewCache(5 * time.Minute),
		pokedex:        make(map[string]pokemonDetails),
		currentPokemon: make(map[string]int),
		commands: map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Display a help message",
				callback:    commandHelp,
			},
			"map": {
				name:        "map",
				description: "Display the names of 20 location areas in the Pokémon world, goes to the next 20 on each consecutive call",
				callback:    commandMap,
			},
			"mapb": {
				name:        "mapb",
				description: "Go back to the previous 20 location areas in the Pokémon world",
				callback:    commandMapb,
			},
			"explore": {
				name:        "explore",
				description: "Explore a location area",
				callback:    commandExplore,
			},
			"catch": {
				name:        "catch",
				description: "Catch a Pokémon in the current area",
				callback:    commandCatch,
			},
			"inspect": {
				name:        "inspect",
				description: "Inspect your Pokedex",
				callback:    commandInspect,
			},
		},
	}

	startRepl(cfg)
}
