package main

import "fmt"

// HangmanTerm is terminal specific structure
type HangmanTerm struct {
}

// GetInput is exported and used in hangman.go for terminal specific inputs
func (h HangmanTerm) GetInput() string {
	// take the input
	str := ""
	fmt.Scanln(&str)
	return str
}

// RenderGame is exported and used in hangman.go for terminal specific renders
func (h HangmanTerm) RenderGame(placeholder []string, entries map[string]bool, chances int) {
	// Console display
	fmt.Println()
	fmt.Println(placeholder) // render the placeholder
	fmt.Println()
	fmt.Println("Chances left: ", chances) // render the chances left
	keys := GetKeys(entries)               // get the wrong guessed keys
	fmt.Println()

	fmt.Println(keys) // show the letters or words guessed till now.
	fmt.Println()
	fmt.Print("Guess a letter or the word: ")
}

func Term() {
	h := HangmanTerm{}

	word := GetWord()

	if result, err := PlayGame(&h, word); result == true {
		fmt.Println("You win! You've saved yourself from a hanging")
	} else {
		fmt.Println(err)
		fmt.Println("Damn! You're hanged!!")
		fmt.Println("Word was: ", word)
	}
}
