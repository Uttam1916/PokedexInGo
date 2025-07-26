package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type config struct {
	Next     string
	Previous string
}

type jsonResult struct {
	Result []locationAreas `json:"results"`
}

type locationAreas struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

// mamke commands globally accesible
var commands = map[string]cliCommand{}

func main() {

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
		description: "Prints 20 locations",
		callback:    commandMap,
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
			fmt.Printf("Unkown command")
			continue

		}

		err := cmd.callback()
		if err != nil {
			fmt.Println("Something went wrong : ", err)
		}
	}
}

func commandHelp() error {
	fmt.Println("Usages:")
	for _, com := range commands {
		fmt.Printf("%v : %v\n", com.name, com.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Thank you for using Pokedex!")
	os.Exit(0)
	return nil
}

func commandMap() error {
	// make the request
	res, err := http.Get("https://pokeapi.co/api/v2/location-area")
	if err != nil {
		return fmt.Errorf("error making request", err)
	}
	defer res.Body.Close()
	// create a slice of bytes from the io reader
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error creating body", err)
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("Status code >299")
	}
	// decode the bytes into a struct
	var areas jsonResult
	err = json.Unmarshal(body, &areas)
	if err != nil {
		return fmt.Errorf("couldnt decode json to struct")
	}
	for _, area := range areas.Result {
		fmt.Printf("%s\n", area.Name)
	}

	return nil
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)

	if text == "" {
		return []string{}
	}
	return strings.Fields(text)
}
