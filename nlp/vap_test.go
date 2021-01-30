package nlp

import (
	"github.com/stretchr/testify/assert"
	"github.com/va-voice-gateway/utils"
	"testing"
)

func TestGetBotConfigs(t *testing.T) {
	botConfigs, err := GetBotConfigs()
	if err != nil {
		t.Fatalf("TestGetBotConfigs failed %v", err)
	}

	utils.PrettyPrint(botConfigs)

	assert.Equal(t, botConfigs[0].Name, "Freight Customer Service Voice", "botConfigs[0].Name does not match")
	assert.Equal(t, botConfigs[0].Channels.VoiceGW.TTSContentUrl, "http://localhost:8444", "botConfigs[0].Channels.VoiceGW.TTSContentUrl does not match")

}
