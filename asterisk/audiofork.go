package asterisk

import (
	"github.com/gorilla/websocket"
	"github.com/va-voice-gateway/gateway/config"
	"github.com/va-voice-gateway/stt/google"
	"github.com/va-voice-gateway/utils"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

// see https://tutorialedge.net/golang/go-websocket-tutorial/
func AudioForkHandler(w http.ResponseWriter, r *http.Request, channelId *string, botId *string, lang *string) {
	log.Printf("AudioForkHandler called for channel: %v botId: %v lang: %v\n", *channelId, *botId, *lang)

	botConfigs := config.BotConfigs(nil)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error when upgrading websocket connection", err)
		return
	}
	defer conn.Close()

	audioStream := make(chan []byte)

	// TBD: GetSTTBotConfig tailored now for google STT only!
	recognitionConfig := botConfigs.GetSTTBotConfig(botId, lang)
	if recognitionConfig == nil {
		log.Printf("Unable to find STT config for %v %v\n", *botId, *lang)
		return
	}

	log.Printf("AudioForkHandler: recognitionConfig: \n")
	utils.PrettyPrint(recognitionConfig)

	// TBD: call here either google or ms stt based on config
	go google.PerformGoogleSTT(audioStream, recognitionConfig, botId, channelId, lang)

	log.Printf("AudioForkHandler: entering loop")

	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error when reading AudioFork message", err)
			return
		} else if mt != websocket.BinaryMessage {
			log.Println("Received wrong AudioFork message type", mt)
		} else {
			audioStream <- p
			continue
		}
	}
	log.Println("AudioForkHandler: loop left")
}
