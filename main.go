package main

import (
	"strings"
)

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)

	if text == "" {
		return []string{}
	}
	return strings.Fields(text)
}
