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

	for {
		guesses := []string{}
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
		fmt.Println()
		fmt.Println(placeholder) // render the placeholder
		fmt.Println()
		fmt.Println(chances) // render the chances left
		for i := range entries {
			guesses = append(guesses, i)
		}
		fmt.Println()

		fmt.Println(guesses) // show the letters or words guessed till now.
		fmt.Println()
		fmt.Print("Guess a letter or the word: ")

		// take the input
		str := ""
		fmt.Scanln(&str)

		// compare and update entries, placeholder and chances.

		_, ok := entries[str]

		if ok {
			continue
		}

		if str == word {
			fmt.Println("You win!!!")
			break
		}

		found := false
		for i, v := range word {
			if string(v) == str {
				//update the placeholder
				placeholder[i] = str
				found = true
			}
		}
		if !found {
			chances -= 1
			entries[str] = true
		}

	}
}
