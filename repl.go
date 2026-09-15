package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

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
	pokedex        map[string]pokemonDetails
	currentPokemon map[string]int
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
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
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
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	fmt.Println()

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

func commandInspect(cfg *config, args ...string) error {
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

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please enter a Pokémon name")
	}

	pokemonName := args[0]

	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	var pokemon pokemonDetails

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

		const maxXP = 300
		chance := maxXP - pokemon.BaseExperience
		roll := rand.Intn(maxXP)

		if roll < chance {
			fmt.Printf("%v caught\n", pokemon.Name)
			cfg.pokedex[pokemon.Name] = pokemon
		} else {
			fmt.Printf("%v escaped", pokemon.Name)
		}
	}
	fmt.Println()
	return nil
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			return
		}
		words := cleanInput(scanner.Text())
		command := words[0]
		args := words[1:]

		c, ok := cfg.commands[command]
		if !ok {
			fmt.Println("Unknown value")
			continue
		}

		err := c.callback(cfg, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}
