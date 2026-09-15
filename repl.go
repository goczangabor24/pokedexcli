package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
	callback    func(*config) error
}

type config struct {
	commands map[string]cliCommand
	next     string
	previous string
	cache    *pokecache.Cache
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

func commandExit(cfg *config) error {
	fmt.Println()
	fmt.Println("Closing the Pokedex... Goodbye!")
	fmt.Println()
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!\nUsage: ")
	fmt.Println()

	for _, command := range cfg.commands {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	fmt.Println()

	return nil
}

func commandMap(cfg *config) error {
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

func commandMapb(cfg *config) error {
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

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			return
		}
		command := cleanInput(scanner.Text())[0]

		c, ok := cfg.commands[command]
		if !ok {
			fmt.Println("Unknown value")
			continue
		}

		err := c.callback(cfg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
