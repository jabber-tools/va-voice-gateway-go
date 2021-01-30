package nlp

import (
	"encoding/json"
	"github.com/va-voice-gateway/gateway/config"
	"io/ioutil"
	"log"
)

// for now just taking from file and parsing
// in the future VAP API will be called:
// /vapapi/vap-mgmt/config-mgmt/v1?voiceEnabled=1
func GetBotConfigs() ([]config.BotConfig, error) {

	var botConfigs [] config.BotConfig

	content, err := ioutil.ReadFile("c:/tmp/botconfigs.json")
	if err != nil {
		log.Printf("GetBotConfigs ReadFile error: %v", err)
		return nil, err
	}

	data := []byte(content)

	err = json.Unmarshal(data, &botConfigs)

	if err != nil {
		log.Printf("GetBotConfigs Unmarshal error: %v", err)
		return nil, err
	}

	return botConfigs, nil
}