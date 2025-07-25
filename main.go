package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)

	if text == "" {
		return []string{}
	}
	return strings.Fields(text)
}
