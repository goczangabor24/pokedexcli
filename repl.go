package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"

	"github.com/goczangabor24/pokedexcli/internal/pokecache"
)

func cleanInput(text string) []string {
	return strings.Split(
		strings.TrimSpace(strings.ToLower(text)),
		" ",
	)
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	commands       map[string]cliCommand
	next           string
	previous       string
	cache          *pokecache.Cache
	pokedex        map[string]*pokemonDetails
	currentPokemon map[string]int
	currentArea    []string
	pokemonStats   map[string]map[string]int
	pokemonToFight map[string]map[string]int
}

type locationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type locationAreaDetailResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type pokemonDetails struct {
	XP             int
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println()
	fmt.Println("Closing the Pokedex... Goodbye!")
	fmt.Println()
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage: ")
	fmt.Println()

	for _, command := range cfg.commands {
		fmt.Printf("- %v: %v\n\n", command.name, command.description)
	}
	fmt.Println()

	return nil
}

func commandCurrent(cfg *config, args ...string) error {
	if len(args) == 0 {

		fmt.Println("\nIf you'd like to see the current area type 'current area',\nif you'd like to see the current Pokemon type 'current pokemon'\n")

	} else if args[0] == "area" {
		if len(cfg.currentArea) == 0 {
			fmt.Println("\nStart exploring the area with the 'map' command\n")
			return nil
		}

		fmt.Println()

		for _, area := range cfg.currentArea {
			fmt.Println(area)
		}
		fmt.Println()
		return nil

	} else if args[0] == "pokemon" {
		if len(cfg.currentArea) == 0 {
			fmt.Println("\nExplore the pokemon in a given area with the 'explore' command\n")
			return nil
		}
		fmt.Println()

		for pokemon, _ := range cfg.currentPokemon {
			fmt.Println(pokemon)
		}

		fmt.Println()
		return nil
	}
	return nil
}

func commandMap(cfg *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area/"

	if cfg.next != "" {
		url = cfg.next
	}

	var locations locationAreaResponse

	cachedRes, ok := cfg.cache.Get(url)
	if ok {
		if err := json.Unmarshal(cachedRes, &locations); err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &locations); err != nil {
			return err
		}
		cfg.cache.Add(url, data)
	}
	cfg.currentArea = []string{}

	fmt.Println()
	for _, area := range locations.Results {
		fmt.Println(area.Name)
		cfg.currentArea = append(cfg.currentArea, area.Name)
	}
	fmt.Println()

	if locations.Next != nil {
		cfg.next = *locations.Next
	} else {
		cfg.next = ""
	}

	if locations.Previous != nil {
		cfg.previous = *locations.Previous
	} else {
		cfg.previous = ""
	}

	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}

	url := cfg.previous

	var locations locationAreaResponse

	cachedRes, ok := cfg.cache.Get(url)
	if ok {
		if err := json.Unmarshal(cachedRes, &locations); err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &locations); err != nil {
			return err
		}
		cfg.cache.Add(url, data)
	}

	fmt.Println()
	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	fmt.Println()

	if locations.Next != nil {
		cfg.next = *locations.Next
	} else {
		cfg.next = ""
	}

	if locations.Previous != nil {
		cfg.previous = *locations.Previous
	} else {
		cfg.previous = ""
	}

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a location area")
	}

	areaName := args[0]

	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	var pokemon locationAreaDetailResponse

	cachedRes, ok := cfg.cache.Get(url)
	if ok {
		if err := json.Unmarshal(cachedRes, &pokemon); err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &pokemon); err != nil {
			return err
		}
		cfg.cache.Add(url, data)
	}
	cfg.currentPokemon = make(map[string]int)

	fmt.Println()
	fmt.Printf("Exploring %s...\n", areaName)
	fmt.Println("Found Pokemon: ")

	for _, encounter := range pokemon.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
		cfg.currentPokemon[encounter.Pokemon.Name] = 1
	}
	fmt.Println()
	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	fmt.Println()
	if len(cfg.pokedex) == 0 {
		fmt.Println("Your pokedex is empty, start catchin' em all!")
	}

	for pokemon, _ := range cfg.pokedex {
		fmt.Println(pokemon)
	}
	fmt.Println()
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		fmt.Println()
		return fmt.Errorf("Please enter a Pokemon to inspect\n")
	}

	fmt.Println()
	if _, ok := cfg.pokedex[args[0]]; ok {
		pokemon := cfg.pokedex[args[0]]

		fmt.Printf("XP: %v\n", pokemon.XP)
		fmt.Printf("Name: %v\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)

		fmt.Println("Stats: ")
		for _, stat := range pokemon.Stats {
			fmt.Printf(" - %v: %d\n", stat.Stat.Name, stat.BaseStat)
		}

		fmt.Println("Types: ")
		for _, pokemonType := range pokemon.Types {
			fmt.Printf(" - %v\n", pokemonType.Type.Name)
		}

		fmt.Println()
		cfg.pokedex[args[0]].XP += 1
	} else {
		fmt.Println("You need to catch this Pokemon first!")
	}
	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please enter a Pokémon name")
	}

	pokemonName := args[0]

	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	var pokemon *pokemonDetails

	cachedRes, ok := cfg.cache.Get(url)
	if ok {
		if err := json.Unmarshal(cachedRes, &pokemon); err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &pokemon); err != nil {
			return err
		}
		cfg.cache.Add(url, data)
	}

	fmt.Println()

	if pokemon.Name == "" {

		fmt.Println("Invalid Pokemon name")

	} else if cfg.currentPokemon[pokemon.Name] != 1 {

		fmt.Println("No such Pokemon found in the area")

	} else if _, ok := cfg.pokedex[pokemon.Name]; ok {

		fmt.Printf("You have already caught %v\n", pokemon.Name)

	} else {

		fmt.Printf("Throwing Pokeball at %s ...\n", pokemon.Name)
		time.Sleep(1 * time.Second)

		const maxXP = 500
		chance := maxXP - pokemon.BaseExperience
		roll := rand.Intn(maxXP)

		if roll < chance {
			fmt.Printf("%v caught\n", pokemon.Name)
			cfg.pokedex[pokemon.Name] = pokemon
			cfg.pokemonStats[pokemon.Name] = make(map[string]int)
			for _, stat := range pokemon.Stats {
				cfg.pokemonStats[pokemon.Name][stat.Stat.Name] = stat.BaseStat
			}
			cfg.pokemonStats[pokemon.Name]["Base Experience"] = pokemon.BaseExperience
		} else {
			fmt.Printf("%v escaped", pokemon.Name)
		}
		fmt.Println()
	}
	return nil
}

func startRepl(cfg *config) {
	fmt.Println("Welcome to the Pokedex, type in the 'help' command to learn how to use it!\nHave fun!")

	rl, err := readline.New("Pokedex > ")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()

		if err == readline.ErrInterrupt {
			continue
		}

		if err == io.EOF {
			return
		}

		words := cleanInput(input)

		if len(words) == 0 {
			continue
		}

		command := words[0]
		args := words[1:]

		c, ok := cfg.commands[command]
		if !ok {
			fmt.Println("Unknown value")
			continue
		}

		err = c.callback(cfg, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}
