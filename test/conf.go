package test

import "fmt"

// Hello returns a greeting for the named person.
func Hello() string {
	var name string
	name = "ifile"
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}
