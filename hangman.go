package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var offline = flag.Bool("offline", false, "elephant")

type hangmanterm struct {
	entries     map[string]bool // lookup for entries made by the user.
	placeholder []string        // list of "_" corrosponding to the number of letters in the word. [ _ _ _ _ _ ]
	word        string
	chances     int
}

func getword() string {
	if *offline { // if dev flag is passed
		return "elephant"
	}

	resp, err := http.Get("https://random-word-api.herokuapp.com/word?number=5") // requestinng random 5 words to an external api.
	if err != nil {
		return "elephant"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "elephant"
	}
	words := []string{}
	err = json.Unmarshal(body, &words) // Unmarshal the json object into words slice
	if err != nil {
		return "elephant"
	}

	for _, v := range words {
		if len(v) > 4 && len(v) < 9 {
			return v
		}
	}
	return "elephant"
}

func getkeys(entries map[string]bool) (keys []string) {
	for i := range entries {
		keys = append(keys, i)
	}
	return
}

func playGame(h1 *hangmanterm, done chan bool) {
	for {
		// evaluate a loss! If user guesses a wrong letter or the wrong word, they lose a chance.
		userInput := strings.Join(h1.placeholder, "")
		if h1.chances == 0 && userInput != h1.word {
			done <- false
			break
		}
		if userInput == h1.word {
			done <- true
			break
		}

		h1.RenderGame()

		// take the input
		str := h1.GetInput()

		// compare and update entries, placeholder and chances.

		if str == h1.word { // if user correct word instead of letter.
			done <- true
			break
		}
		_, ok := h1.entries[str] // check for duplicated guess.

		if ok {
			continue // if duplicated do nothing and continue.
		}

		found := false
		for i, v := range h1.word {
			if string(v) == str { // check the presence of guessed key in word
				//update the placeholder
				h1.placeholder[i] = str
				found = true
			}
		}
		if !found {
			h1.chances -= 1
			h1.entries[str] = true
		}
	}
}

func main() {
	flag.Parse()
	h1 := hangmanterm{
		placeholder: []string{},
		entries:     map[string]bool{},
		chances:     8,
		word:        "elephant",
	}

	h1.word = getword()
	for i := 0; i < len(h1.word); i++ {
		h1.placeholder = append(h1.placeholder, "_")
	}

	ticker := time.NewTicker(30 * time.Second) // a channel that will deliver time after every 30sec
	defer ticker.Stop()                        // stops the ticker
	done := make(chan bool)
	go playGame(&h1, done) // stared a new go routine as playGame.

	select { // select either one of the data from the channel.
	case r := <-done: // if user guessed correct word.
		if r {
			fmt.Println("Excellent, You win!!")
		} else {
			fmt.Println("You loss! Try Again")
			fmt.Println("Word: ", h1.word)
		}
		break
	case <-ticker.C: // if the time crosses 30secs(1st tick).
		fmt.Println("\nTime up!!!")
		fmt.Println("Word: ", h1.word)
		break
	}
}
