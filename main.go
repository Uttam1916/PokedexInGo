package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"math/rand"

	"github.com/Uttam1916/PokedexInGo/internal/pokecache"
)

var cache *pokecache.Cache
var user_Pokedex map[string]pokepoke

type config struct {
	next     string
	previous string
}

type jsonSpecificlocationArea struct {
	Location locationAreas      `json:"location"`
	Pokemons []pokemonEncounter `json:"pokemon_encounters"`
}

type jsonLocationArea struct {
	Result   []locationAreas `json:"results"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
}

type locationAreas struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type pokemonEncounter struct {
	Pokemon pokemon `json:"pokemon"`
}

type pokemon struct {
	PokemonName string `json:"name"`
	PokemonUrl  string `json:"url"`
}
type pokepoke struct {
	Base_experience int    `json:"base_experience"`
	Pokepoke        string `json:"name"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
	callbackwp  func(c *config, url string) error
}

// mamke commands globally accesible
var commands = map[string]cliCommand{}

func main() {
	// initializing cache and url
	cache = pokecache.NewCache(1 * time.Minute)

	c := config{
		next: "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
	}

	user_Pokedex = make(map[string]pokepoke)

	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exits the program",
		callback:    commandExit,
	}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
	commands["map"] = cliCommand{
		name:        "map",
		description: "Prints next 20 locations",
		callback:    commandMap,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Prints previous 20 locations",
		callback:    commandMapb,
	}
	commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "lists caught pokemon",
		callback:    showPokemon,
	}
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "explores area for pokemon",
		callbackwp:  commandExplore,
	}
	commands["catch"] = cliCommand{
		name:        "catch",
		description: "tries to catch a pokemon",
		callbackwp:  commandCatch,
	}

	// create a scanner to read line by line
	scanner := bufio.NewScanner(os.Stdin)
	//infinite for loop to wait for user input
	fmt.Println("Welcome to pokedex!")

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())

		//check if command exists
		cmd, ok := commands[input[0]]
		if !ok {
			fmt.Printf("Unkown command\n")
			continue

		}
		if len(input) == 1 {
			err := cmd.callback(&c)
			if err != nil && input[1] != "" {
				fmt.Println("Something went wrong : ", err)
			}
		} else {
			err := cmd.callbackwp(&c, input[1])
			if err != nil {
				fmt.Println("Something went wrong : ", err)
			}
		}
	}
}

func commandHelp(c *config) error {
	fmt.Println("Usages:")
	for _, com := range commands {
		fmt.Printf("%v : %v\n", com.name, com.description)
	}
	return nil
}

func commandExit(c *config) error {
	fmt.Println("Thank you for using Pokedex!")
	os.Exit(0)
	return nil
}

func commandMapb(c *config) error {
	if c.previous == "" {
		fmt.Println("already at the beginning")
		return nil
	}
	err := fetchdata(c.previous, c)
	if err != nil {
		fmt.Println("error fetching data")
	}
	return nil
}

func commandMap(c *config) error {
	if c.next == "" {
		fmt.Println("already at the end")
		return nil
	}
	err := fetchdata(c.next, c)
	if err != nil {
		fmt.Println("error fetching data")
	}
	return nil
}

func fetchdata(url string, c *config) error {
	//try to get from cache
	if data, ok := cache.Get(url); ok {
		return processData(data, c)
	}

	// make the request
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making request", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("Status code >299")
	}
	// create a slice of bytes from the io reader
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error creating body", err)
	}
	cache.Add(url, body)
	return processData(body, c)

}

func processData(data []byte, c *config) error {
	var areas jsonLocationArea
	err := json.Unmarshal(data, &areas)
	if err != nil {
		return fmt.Errorf("couldnt decode json to struct")
	}
	for _, area := range areas.Result {
		fmt.Printf("%s\n", area.Name)
	}
	//change the config to paginate
	c.next = areas.Next
	c.previous = areas.Previous
	return nil
}

func commandExplore(c *config, location_name string) error {
	res, err := http.Get("https://pokeapi.co/api/v2/location-area/" + location_name)
	if err != nil {
		fmt.Println("error requesting poke-encounters")
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading body")
		return nil
	}
	//print pokemons present in area
	var specific_area jsonSpecificlocationArea
	err = json.Unmarshal(body, &specific_area)
	for _, pokemon := range specific_area.Pokemons {
		fmt.Println(pokemon.Pokemon.PokemonName)
	}
	return nil
}

func commandCatch(c *config, pokemon_name string) error {
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + pokemon_name)
	if err != nil {
		fmt.Println("error requesting pokemon-stats")
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading body")
		return nil
	}
	var catching_poke pokepoke
	err = json.Unmarshal(body, &catching_poke)
	// try catching the pokemon
	fmt.Printf("Throwing a pokeball at %s...\n", pokemon_name)
	println("Trying...")
	time.Sleep(time.Millisecond * 700)
	println("Trying...")
	time.Sleep(time.Millisecond * 700)
	println("Almost there..")
	time.Sleep(time.Millisecond * 1000)
	// Calculate catch chance
	// baseExperience: ~50 (easy) to 300+ (hard)
	catchChance := 1.0 / (float64(catching_poke.Base_experience)/50.0 + 1.0)
	if catchChance > 0.8 {
		catchChance = 0.8
	}
	if rand.Float64() < catchChance {
		fmt.Printf("You caught %s!\n", pokemon_name)
		user_Pokedex[pokemon_name] = catching_poke
	} else {
		fmt.Printf("%s escaped!\n", pokemon_name)
	}

	return nil
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.TrimSpace(text))
	return words
}

func showPokemon(c *config) error {
	i := 1
	for name, pokestruct := range user_Pokedex {
		fmt.Printf("%v. %v XP: %d \n", i, name, pokestruct.Base_experience)
		i++
	}
	return nil
}
