package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/smcallister/pokedex/internal/pokeapi"
)

type cmdConfig struct {
	next 		*string
	previous	*string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*cmdConfig, ...string) error
}

var commands map[string]cliCommand
var pokedex map[string]*pokeapi.Pokemon

func cleanInput(text string) []string {
	var cleaned []string

	// Break up the input text into words, trimming spaces.
	words := strings.Fields(text)
	for _, word := range words {
		// Convert to lowercase and trim spaces.
		cleanedWord := strings.TrimSpace(word)
		cleaned = append(cleaned, strings.ToLower(cleanedWord))
	}

	return cleaned
}

func commandHelp(config *cmdConfig, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(config *cmdConfig, args ...string) error {
	// Get the next page of location areas.
	if config.next == nil {
		fmt.Println("you're on the last page")
	}

	page, err := pokeapi.GetLocationAreas(*config.next)
	if err != nil {
		return err
	}

	// Print the names of the location areas.
	for _, result := range page.Results {
		fmt.Println(result.Name)
	}

	config.next = page.Next
	config.previous = page.Previous
	return nil
}

func commandMapB(config *cmdConfig, args ...string) error {
	// Get the previous page of location areas.
	if config.previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	page, err := pokeapi.GetLocationAreas(*config.previous)
	if err != nil {
		return err
	}

	// Print the names of the location areas.
	for _, result := range page.Results {
		fmt.Println(result.Name)
	}

	// Update the config.
	config.next = page.Next
	config.previous = page.Previous
	return nil
}

func commandExplore(config *cmdConfig, args ...string) error {
	// Check if a location area is specified.
	if len(args) == 0 {
		return fmt.Errorf("Please specify a location area to explore")
	}

	// Get the location area.
	fmt.Printf("Exploring %s\n", args[0])
	area, err := pokeapi.GetLocationArea(args[0])
	if err != nil {
		return err
	}

	// Print the details of the location area.
	fmt.Println("Found Pokemon:");
	for _, encounter := range area.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *cmdConfig, args ...string) error {
	// Check if a pokemon is specified.
	if len(args) == 0 {
		return fmt.Errorf("Please specify a pokemon to catch")
	}

	// Get the pokemon.
	pokemon, err := pokeapi.GetPokemon(args[0])
	if err != nil {
		return err
	}

	// Try to catch the pokemon.
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if rand.Intn(100000) > pokemon.BaseExperience {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		fmt.Println("You may now inspect it with the inspect command.")
		pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(config *cmdConfig, args ...string) error {
	// Check if a pokemon is specified.
	if len(args) == 0 {
		return fmt.Errorf("Please specify a pokemon to inspect")
	}

	// Get the pokemon from the Pokedex.
	pokemon, exists := pokedex[args[0]]
	if !exists {
		return fmt.Errorf("you have not caught that pokemon")
	}

	// Print the details of the pokemon.
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("- %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("- %s\n", typeInfo.Type.Name)
	}

	return nil
}

func commandPokedex(config *cmdConfig, args ...string) error {
	if len(pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Catch some Pokemon first!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Printf("- %s\n", name)
	}

	return nil
}

func commandExit(config *cmdConfig, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!");
	os.Exit(0)
	return nil
}

func main() {
	// Initialize the commands.
	commands = map[string]cliCommand{
    	"help": {
        	name:        "help",
        	description: "Displays a help message",
        	callback:    commandHelp,
    	},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map",
			description: "Displays the names of the previous 20 location areas in the Pokemon world",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore <area name>",
			description: "Explore a location area in the Pokemon world",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon name>",
			description: "Catch a pokemon and add it to your Pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon name>",
			description: "Inspects a pokemon in your Pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Prints the names of all pokemon in your Pokedex",
			callback:    commandPokedex,
		},
    	"exit": {
        	name:        "exit",
        	description: "Exit the Pokedex",
        	callback:    commandExit,
    	},
	}

	pokedex = make(map[string]*pokeapi.Pokemon)

	url := pokeapi.LocationAreaURL
	config := cmdConfig{
		next:     &url,
		previous: nil,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Prompt the user for input.
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break // Exit on EOF
		}

		// Read the input and clean it.
		input := scanner.Text()
		cleanedWords := cleanInput(input)
		if len(cleanedWords) == 0 {
			break
		}

		// Look up the command and execute it.
		if command, exists := commands[cleanedWords[0]]; exists {
			if err := command.callback(&config, cleanedWords[1:]...); err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", cleanedWords[0])
		}
	}
}
