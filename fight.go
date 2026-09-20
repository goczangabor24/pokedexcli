package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func fight(cfg *config, args ...string) error {
	if len(args) < 2 {
		fmt.Println("\nFirst enter the Pokemon you want to fight, then one of your own Pokemon from your Pokedex")
		return nil
	}

	enemyName := args[0]
	playerName := args[1]

	if !validFight(cfg, enemyName, playerName) {
		return nil
	}

	cfg.cooldownsMu.Lock()
	onCooldown := cfg.cooldowns[playerName]
	cfg.cooldownsMu.Unlock()

	if onCooldown {
		fmt.Printf("\n%s is recovering and can't fight yet!\n", playerName)
		return nil
	}

	enemy, err := getPokemonForFight(cfg, enemyName)
	if err != nil {
		return err
	}

	addEnemyStats(cfg, enemy)

	playerWon := runBattle(cfg, enemyName, playerName)

	if playerWon {
		catchPokemonAfterFight(cfg, enemy)
	}

	levelUpPokemon(cfg, playerName)

	delete(cfg.pokemonToFight, enemyName)

	startCooldown(cfg, playerName)

	return nil
}

func validFight(cfg *config, enemyName, playerName string) bool {
	if cfg.currentPokemon[enemyName] != 1 {
		fmt.Println("No such Pokemon in the area.\nCheck the available Pokemon with 'current pokemon'")
		return false
	}

	if _, ok := cfg.pokedex[playerName]; !ok {
		fmt.Println("Pick a Pokemon you've already caught.\nHint: check you Pokedex")
		return false
	}

	if _, ok := cfg.pokedex[enemyName]; ok {
		fmt.Println("You've already caught this Pokemon")
		return false
	}

	return true
}

func getPokemonForFight(cfg *config, name string) (*pokemonDetails, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	var pokemon *pokemonDetails

	cachedRes, ok := cfg.cache.Get(url)
	if ok {
		if err := json.Unmarshal(cachedRes, &pokemon); err != nil {
			return nil, err
		}
		return pokemon, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &pokemon); err != nil {
		return nil, err
	}

	cfg.cache.Add(url, data)

	return pokemon, nil
}

func addEnemyStats(cfg *config, pokemon *pokemonDetails) {
	cfg.pokemonToFight[pokemon.Name] = make(map[string]int)

	for _, stat := range pokemon.Stats {
		cfg.pokemonToFight[pokemon.Name][stat.Stat.Name] = stat.BaseStat
	}
}

func runBattle(cfg *config, enemyName, playerName string) bool {
	enemyHp := cfg.pokemonToFight[enemyName]["hp"]
	enemyAttack := cfg.pokemonToFight[enemyName]["attack"]
	enemyDefense := cfg.pokemonToFight[enemyName]["defense"]

	playerHp := cfg.pokemonStats[playerName]["hp"]
	playerAttack := cfg.pokemonStats[playerName]["attack"]
	playerDefense := cfg.pokemonStats[playerName]["defense"]

	for enemyHp > 0 && playerHp > 0 {

		playerDamage := playerAttack - enemyDefense
		if playerDamage < 1 {
			playerDamage = 1
		}

		enemyHp -= playerDamage
		if enemyHp <= 0 {
			printFightStatus(enemyHp, playerHp)
			fmt.Println("Player won")
			return true
		}

		enemyDamage := enemyAttack - playerDefense
		if enemyDamage < 1 {
			enemyDamage = 1
		}

		playerHp -= enemyDamage

		printFightStatus(enemyHp, playerHp)
	}
	fmt.Println("\nCPU won")
	return false

}

func printFightStatus(enemyHp, playerHp int) {
	if enemyHp < 0 {
		enemyHp = 0
	}

	if playerHp < 0 {
		playerHp = 0
	}

	fmt.Printf("\nEnemy: %v\n", enemyHp)
	fmt.Printf("Player: %v\n", playerHp)

	time.Sleep(500 * time.Millisecond)
}

func catchPokemonAfterFight(cfg *config, pokemon *pokemonDetails) {
	fmt.Printf("\n%v caught!", pokemon.Name)

	cfg.pokedex[pokemon.Name] = pokemon
	cfg.pokemonStats[pokemon.Name] = make(map[string]int)

	for _, stat := range pokemon.Stats {
		cfg.pokemonStats[pokemon.Name][stat.Stat.Name] = stat.BaseStat
	}

	cfg.pokemonStats[pokemon.Name]["Base Experience"] = pokemon.BaseExperience
}

func levelUpPokemon(cfg *config, pokemonName string) {

	for stat := range cfg.pokemonStats[pokemonName] {
		if stat == "Base Experience" {
			continue
		}

		cfg.pokemonStats[pokemonName][stat]++
	}

	pokemon := cfg.pokedex[pokemonName]

	for i := range pokemon.Stats {
		pokemon.Stats[i].BaseStat++
	}
}

func startCooldown(cfg *config, pokemonName string) {
	cfg.cooldownsMu.Lock()
	cfg.cooldowns[pokemonName] = true
	cfg.cooldownsMu.Unlock()

	go func() {
		time.Sleep(15 * time.Second)

		cfg.cooldownsMu.Lock()
		cfg.cooldowns[pokemonName] = false
		cfg.cooldownsMu.Unlock()

		fmt.Printf("\n%s has recovered and can fight again!\n", pokemonName)
	}()
}
