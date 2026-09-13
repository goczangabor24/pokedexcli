package main

import "strings"

func cleanInput(text string) []string {
	result := strings.Split(strings.TrimSpace(strings.ToLower(text)), " ")
	return result
}
