package main

import (
	"time"

	"github.com/goczangabor24/pokedexcli/internal/pokecache"
)

func main() {
	cfg := &config{
		cache:          pokecache.NewCache(5 * time.Minute),
		pokedex:        make(map[string]*pokemonDetails),
		currentPokemon: make(map[string]int),
		currentArea:    []string{},
		pokemonStats:   make(map[string]map[string]int),
		pokemonToFight: make(map[string]map[string]int),
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
			"pokedex": {
				name:        "pokedex",
				description: "Display your Pokedex entries",
				callback:    commandPokedex,
			},
			"inspect": {
				name:        "inspect",
				description: "Inspect the details of a Pokemon you caught",
				callback:    commandInspect,
			},
			"current": {
				name:        "current",
				description: "current area: Shows current areas; current pokemon: Shows current Pokemon",
				callback:    commandCurrent,
			},
			"fight": {
				name:        "fight",
				description: "After the 'fight' type a Pokemon from the area you wish to catch then type a second Pokemon from your Pokedex to fight it with",
				callback:    fight,
			},
		},
	}

	startRepl(cfg)
}
