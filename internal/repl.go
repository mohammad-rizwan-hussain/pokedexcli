package internal

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/mrizwan/pokedexcli/internal/pokeapi"
)

type config struct {
	pokeapiClient pokeapi.Client
	commands      map[string]cliCommand
	pokedex       *Pokedex
}

func NewConfig(timeout, cacheDuration time.Duration) *config {
	client := pokeapi.NewClient(timeout, cacheDuration) // 10 seconds timeout, 5 minutes cache duration
	cfg := &config{
		pokeapiClient: client,
		commands:      getCommands(),
		pokedex:       NewPokedex(),
	}
	return cfg
}

func StartRepl(cfg *config) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if ok := reader.Scan(); !ok {
			fmt.Printf("Error reading input: %w\n Exiting...", reader.Err())
			break
		}

		words := CleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}
		commandName := words[0]

		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg, words[1:])
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func CleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas in the Pokemon world",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays the names of 10 Pokemons from the area in the Pokemon world",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon from the area in the Pokemon world",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon in your Pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays the names of all Pokemons in your Pokedex",
			callback:    commandPokedex,
		},
	}
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	cmds := ""
	for _, cmd := range cfg.commands {
		cmds += fmt.Sprintf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Printf(`Welcome to the Pokedex!
Usage:

%s
`, cmds)
	return nil
}

func commandMap(cfg *config, args []string) error {
	areas, err := cfg.pokeapiClient.GetLocationAreas(false)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %v", err)
	}
	for _, area := range areas {
		fmt.Println(area)
	}

	return nil
}

func commandMapb(cfg *config, args []string) error {
	areas, err := cfg.pokeapiClient.GetLocationAreas(true)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %v", err)
	}
	for _, area := range areas {
		fmt.Println(area)
	}

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}
	areaName := args[0]

	fmt.Printf("Exploring %s...\n", areaName)
	pokemons, err := cfg.pokeapiClient.GetPokemonsInArea(areaName)
	if err != nil {
		return fmt.Errorf("failed to fetch pokemons in area %s: %v", areaName, err)
	}

	for _, pokemon := range pokemons {
		fmt.Printf(" - %s\n", pokemon)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}
	pokemonName := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	pokemon, err := cfg.pokeapiClient.GetPokemonDetails(pokemonName)
	if err != nil {
		return fmt.Errorf("failed to fetch details for pokemon %s: %v", pokemonName, err)
	}
	// typeNames := make([]string, 0, len(pokemon.Types))
	// for _, t := range pokemon.Types {
	// 	typeNames = append(typeNames, t.Type.Name)
	// }
	// stats := make(map[string]int, len(pokemon.Stats))
	// for _, s := range pokemon.Stats {
	// 	stats[s.Stat.Name] = s.BaseStat
	// }
	attemptCatch := func(baseExperience int) bool {
		chance := 100.0 / (1.0 + float64(baseExperience)/50.0)

		if chance < 5 {
			chance = 5
		}

		return rand.Float64()*100 < chance
	}

	if catch := attemptCatch(pokemon.BaseExperience); catch {
		cfg.pokedex.Add(pokemon)
		fmt.Printf("%s was caught!\n", pokemon.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	// fmt.Printf("You caught a %s!\n", pokemon.Name)
	// fmt.Printf("Type: %s\n", strings.Join(typeNames, ", "))
	// fmt.Printf("Stats: %v\n", stats)

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}
	pokemonName := args[0]

	pokemon, exists := cfg.pokedex.Get(pokemonName)
	if !exists {
		return fmt.Errorf("you have not cayght that pokemon")
	}

	typeNames := make([]string, 0, len(pokemon.Types))
	for _, t := range pokemon.Types {
		typeNames = append(typeNames, t.Type.Name)
	}
	stats := make(map[string]int, len(pokemon.Stats))
	for _, s := range pokemon.Stats {
		stats[s.Stat.Name] = s.BaseStat
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for statName, statValue := range stats {
		fmt.Printf("  -%s: %d\n", statName, statValue)
	}
	fmt.Println("Types:")
	for _, typeName := range typeNames {
		fmt.Printf("  - %s\n", typeName)
	}

	return nil
}

func commandPokedex(cfg *config, args []string) error {
	if len(*cfg.pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Catch some Pokemons first!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range *cfg.pokedex {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}
