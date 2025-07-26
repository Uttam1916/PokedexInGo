package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Uttam1916/PokedexInGo/internal/pokecache"
)

var cache *pokecache.Cache

type config struct {
	next     string
	previous string
}

type jsonResult struct {
	Result   []locationAreas `json:"results"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
}

type locationAreas struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

// mamke commands globally accesible
var commands = map[string]cliCommand{}

func main() {
	cache = pokecache.NewCache(1 * time.Minute)

	c := config{
		next: "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
	}

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

	// create a scanner to read line by line
	scanner := bufio.NewScanner(os.Stdin)
	//infinite for loop to wait for user input
	fmt.Println("Pokedex > Hello!")

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		//check if command exists
		cmd, ok := commands[input]
		if !ok {
			fmt.Printf("Unkown command\n")
			continue

		}

		err := cmd.callback(&c)
		if err != nil {
			fmt.Println("Something went wrong : ", err)
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
	var areas jsonResult
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
