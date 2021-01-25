package asterisk

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

// see https://tutorialedge.net/golang/go-websocket-tutorial/
func AudioForkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AudioForkHandler called")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error when upgrading websocket connection", err)
		return
	}

	fmt.Println("AudioForkHandler: entering loop")
	defer conn.Close()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error when reading AudioFork message", err)
			return
		} else if messageType != websocket.BinaryMessage {
			fmt.Println("Received wrong AudioFork message type", messageType)
			continue
		} else {
			fmt.Println("AudioFork bytes", p)
		}
	}
	fmt.Println("AudioForkHandler: loop left")
}
