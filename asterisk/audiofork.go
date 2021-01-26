package asterisk

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/va-voice-gateway/appconfig"
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

	//
	// Google Speech To Text - quick & dirty for now
	// TBD: implement google streaming here as per:
	// https://cloud.google.com/speech-to-text/docs/streaming-recognize
	//


	fmt.Println("AudioForkHandler: entering loop")

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
