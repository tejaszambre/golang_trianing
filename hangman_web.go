package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
)

// HangmanWeb is a web specific struct for hangman
type HangmanWeb struct {
	conn *websocket.Conn
}

var upgrader = websocket.Upgrader{} // use default options

// RenderGame is a hangmanweb pointer receiver for rendering
func (h *HangmanWeb) RenderGame(placeholder []string, entries map[string]bool, chances int) {
	c := h.GetDisplayConn().(*websocket.Conn)
	str := fmt.Sprintf("%v Chances left: %d   Guesses: %v  ", placeholder, chances, GetKeys(entries))
	c.WriteMessage(websocket.TextMessage, []byte(str))
}

// GetDisplayConn is a hangmanweb pointer receiver for getting conn
func (h *HangmanWeb) GetDisplayConn() interface{} {
	return h.conn
}

// GetInput is a hangmanweb pointer receiver for getting the inputs
func (h *HangmanWeb) GetInput() string {
	c := h.GetDisplayConn().(*websocket.Conn)

	_, message, err := c.ReadMessage()
	if err != nil {
		fmt.Println("read:", err)
		return ""
	}

	// validate mt is a TextMessage
	return string(message)
}

// InitRouter is a function
func InitRouter() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/new", newGameHandler).Methods(http.MethodGet)
	router.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	return
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/new")
}

func newGameHandler(rw http.ResponseWriter, req *http.Request) {
	h := HangmanWeb{}
	word := GetWord()

	c, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	h.conn = c

	if result, _ := PlayGame(&h, word); result == true {
		c.WriteMessage(websocket.TextMessage, []byte("You win! You've saved yourself from a hanging"))
	} else {
		c.WriteMessage(websocket.TextMessage, []byte("Damn! You're hanged!!"))
		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Word was: %s", word)))
	}
}

// Web is a method which starts the server
func Web() {
	router := InitRouter()
	server := negroni.Classic()
	server.UseHandler(router)

	server.Run(":3000")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };
    document.getElementById("start").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
        }
        ws.onclose = function(evt) {
            ws = null;
        }
        ws.onmessage = function(evt) {
            print(evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("play").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<form>
<button id="start">Start New Game</button>
<p><input id="input" type="text">
<button id="play">Play</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
