package main

import (
	"fmt"
	"strings"
)

func main() {
	word := "elephant"

	// lookup for entries made by the user.
	entries := map[string]bool{}

	// list of "_" corrosponding to the number of letters in the word. [ _ _ _ _ _ ]
	placeholder := []string{}

	for i := 0; i < len(word); i++ {
		placeholder = append(placeholder, "_")
	}

	chances := 8
	guesses := []string{}

	for {
		// evaluate a loss! If user guesses a wrong letter or the wrong word, they lose a chance.
		userInput := strings.Join(placeholder, "")
		if chances == 0 && userInput != word {
			fmt.Println("You loss! Try Again")
			break
		}
		if chances == 0 && userInput == word {
			fmt.Println("You win!!")
			break
		}
		// Console display
		fmt.Println("\n")
		fmt.Println(placeholder) // render the placeholder
		fmt.Println(chances)     // render the chances left
		for i, _ := range entries {
			guesses = append(guesses, i)
		}

		fmt.Println(guesses) // show the letters or words guessed till now.
		fmt.Println("Guess a letter or the word: ")

		// take the input
		str := ""
		fmt.Scanln(&str)

		// compare and update entries, placeholder and chances.
	}
}
