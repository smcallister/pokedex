package main

import (
	"bufio"
	"fmt"
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
	callback    func(*cmdConfig) error
}

var commands map[string]cliCommand

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

func commandHelp(config *cmdConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(config *cmdConfig) error {
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

func commandMapB(config *cmdConfig) error {
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

func commandExit(config *cmdConfig) error {
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
    	"exit": {
        	name:        "exit",
        	description: "Exit the Pokedex",
        	callback:    commandExit,
    	},
	}

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
			if err := command.callback(&config); err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", cleanedWords[0])
		}
	}
}
