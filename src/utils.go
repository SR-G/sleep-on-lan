package main

import (
	"strings"
)

func ReverseMacAddress(address string) string {
	tokens := strings.Split(address, ":")
	last := len(tokens) - 1
	for i := 0; i < len(tokens)/2; i++ {
		tokens[i], tokens[last-i] = tokens[last-i], tokens[i]
	}
	return strings.Join(tokens, ":")
}
