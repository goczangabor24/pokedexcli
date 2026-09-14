package main

func main() {
	cfg := &config{
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
				name:        "pokemap",
				description: "Displays the names of 20 location areas in the Pokemon world, goes to the next 20 on each consecutive call",
				callback:    commandMap,
			},
			"mapb": {
				name:        "pokemapb",
				description: "Goes back to the previous 20 location areas in the Pokemon world",
				callback:    commandMapb,
			},
		},
	}

	startRepl(cfg)
}
