package main

import "fmt"

func (h hangmanterm) GetInput() string {
	// take the input
	str := ""
	fmt.Scanln(&str)
	return str
}

func (h hangmanterm) RenderGame() {
	// Console display
	fmt.Println()
	fmt.Println(h.placeholder) // render the placeholder
	fmt.Println()
	fmt.Println("Chances left: ", h.chances) // render the chances left
	keys := getkeys(h.entries)               // get the wrong guessed keys
	fmt.Println()

	fmt.Println(keys) // show the letters or words guessed till now.
	fmt.Println()
	fmt.Print("Guess a letter or the word: ")
}
