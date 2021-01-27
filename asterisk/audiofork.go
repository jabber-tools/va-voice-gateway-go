package asterisk

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/stt"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

// see https://tutorialedge.net/golang/go-websocket-tutorial/
func AudioForkHandler(w http.ResponseWriter, r *http.Request, appConfig *appconfig.AppConfig) {
	fmt.Println("AudioForkHandler called")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error when upgrading websocket connection", err)
		return
	}
	defer conn.Close()

	audioStream := make(chan []byte)

	// TBD: call here either google or ms stt based on config
	go stt.PerformGoogleSTT(appConfig, audioStream)

	log.Printf("AudioForkHandler: entering loop")

	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error when reading AudioFork message", err)
			return
		} else if mt != websocket.BinaryMessage {
			fmt.Println("Received wrong AudioFork message type", mt)
		} else {
			audioStream <- p
			continue
		}
	}
	fmt.Println("AudioForkHandler: loop left")
}
