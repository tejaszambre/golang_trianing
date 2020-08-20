package main

import "fmt"

// HangmanTerm is terminal specific structure
type HangmanTerm struct {
	entries     map[string]bool // lookup for entries made by the user.
	placeholder []string        // list of "_" corrosponding to the number of letters in the word. [ _ _ _ _ _ ]
	word        string
	chances     int
}

// GetInput is exported and used in hangman.go for terminal specific inputs
func (h HangmanTerm) GetInput() string {
	// take the input
	str := ""
	fmt.Scanln(&str)
	return str
}

// RenderGame is exported and used in hangman.go for terminal specific renders
func (h HangmanTerm) RenderGame() {
	// Console display
	fmt.Println()
	fmt.Println(h.placeholder) // render the placeholder
	fmt.Println()
	fmt.Println("Chances left: ", h.chances) // render the chances left
	keys := GetKeys(h.entries)               // get the wrong guessed keys
	fmt.Println()

	fmt.Println(keys) // show the letters or words guessed till now.
	fmt.Println()
	fmt.Print("Guess a letter or the word: ")
}
