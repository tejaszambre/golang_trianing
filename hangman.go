package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var MAX_CHANCES int = 8

type Hangman interface {
	RenderGame([]string, map[string]bool, int)
	GetInput() string
}

var offline = flag.Bool("offline", false, "elephant")

func GetWord() string {
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

func GetKeys(entries map[string]bool) (keys []string) {
	for i := range entries {
		keys = append(keys, i)
	}
	return
}

func PlayGame(h Hangman, word string) (r bool, err error) {
	placeholder := []string{}
	entries := map[string]bool{}
	chances := MAX_CHANCES

	for i := 0; i < len(word); i++ {
		placeholder = append(placeholder, "_")
	}

	ticker := time.NewTicker(2 * time.Minute) // a channel that will deliver time after every 30sec
	defer ticker.Stop()                       // stops the ticker
	done := make(chan bool)

	go func() {
		for {
			// evaluate a loss! If user guesses a wrong letter or the wrong word, they lose a chance.
			userInput := strings.Join(placeholder, "")
			if chances == 0 && userInput != word {
				done <- false
				break
			}
			if userInput == word {
				done <- true
				break
			}

			h.RenderGame(placeholder, entries, chances)

			// take the input
			str := h.GetInput()

			// compare and update entries, placeholder and chances.
			if str == word { // if user correct word instead of letter.
				done <- true
				break
			}
			_, ok := entries[str] // check for duplicated guess.

			if ok {
				continue // if duplicated do nothing and continue.
			}

			found := false
			for i, v := range word {
				if string(v) == str { // check the presence of guessed key in word
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
	}()

	select { // select either one of the data from the channel.
	case r = <-done: // if user guessed correct word.
		if r {
			err = nil
			return true, err
		} else {
			err = errors.New("You loss! Try Again")
			return false, err
		}
	case <-ticker.C: // if the time crosses 30secs(1st tick).
		err = errors.New("Time up! Try Again")
		return false, err
	}

}

func main() {
	go Web()
	Term()
}
