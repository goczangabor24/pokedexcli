package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func fight(cfg *config, args ...string) error {

	if len(args) == 0 || len(args) == 1 {
		fmt.Println("\nFirst enter the Pokemon you want to fight, then one of your own Pokemon from your Pokedex")
		return nil
	}

	p1 := args[0]
	p2 := args[1]

	if cfg.currentPokemon[p1] != 1 {
		fmt.Println("No such Pokemon in the area.\nCheck the available Pokemon with 'current pokemon'")
		return nil
	}

	_, ok := cfg.pokedex[p2]
	if !ok {
		fmt.Println("Pick a Pokemon you've already caught.\nHint: check your Pokedex")
		return nil
	}

	if p1 == p2 {
		fmt.Println("You've already caught this Pokemon")
		return nil
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + args[0]

	var pokemon1 *pokemonDetails

	cachedRes, ok := cfg.cache.Get(url)
	if ok {
		if err := json.Unmarshal(cachedRes, &pokemon1); err != nil {
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

		if err := json.Unmarshal(data, &pokemon1); err != nil {
			return err
		}
		cfg.cache.Add(url, data)
	}

	cfg.pokemonToFight[pokemon1.Name] = make(map[string]int)
	for _, stat := range pokemon1.Stats {
		cfg.pokemonToFight[pokemon1.Name][stat.Stat.Name] = stat.BaseStat
	}

	enemyHp := cfg.pokemonToFight[p1]["hp"]
	playerHp := cfg.pokemonStats[p2]["hp"]
	enemyAttack := cfg.pokemonToFight[p1]["attack"]
	playerAttack := cfg.pokemonStats[p2]["attack"]
	enemyDefense := cfg.pokemonToFight[p1]["defense"]
	playerDefense := cfg.pokemonStats[p2]["defense"]

	maxHealthp2 := cfg.pokemonStats[p2]["hp"]

	status := func() {
		fmt.Println()
		if enemyHp < 0 {
			fmt.Printf("\nEnemy: 0\n")
		} else {
			fmt.Printf("\nEnemy: %v\n", enemyHp)
		}

		if playerHp < 0 {
			fmt.Printf("\nPlayer: 0\n")
		} else {
			fmt.Printf("Player: %v\n", playerHp)
		}
		time.Sleep(500 * time.Millisecond)
	}

	for enemyHp > 0 && playerHp > 0 {

		if enemyAttack < playerDefense && playerAttack < enemyDefense {
			fmt.Println("Draw")
			return nil
		}

		if playerAttack-enemyDefense < 0 {
			status()
		} else {
			enemyHp -= playerAttack - enemyDefense
			status()
			if enemyHp <= 0 {
				break
			}
		}

		if enemyAttack-playerDefense < 0 {
			status()
		} else {
			playerHp -= enemyAttack - playerDefense
			status()
			if playerHp <= 0 {
				break
			}
		}
	}

	if enemyHp <= 0 {
		fmt.Println("Player won")
	} else {
		fmt.Println("\nCPU won")
	}
	fmt.Println()

	delete(cfg.pokemonToFight, p1)
	playerHp = maxHealthp2
	for _, stat := range cfg.pokemonStats[p2] {
		stat += 1
	}

	return nil
}
