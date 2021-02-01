package nlp

import (
	"github.com/stretchr/testify/assert"
	"github.com/va-voice-gateway/utils"
	"testing"
	"time"
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

func TestIsExpired1(t *testing.T) {
	TokenCreatedTime := time.Now()
	time.Sleep(time.Second * 5)
	TimeNow := time.Now()
	assert.Equal(t, true, IsExpired(TimeNow, TokenCreatedTime, 3), "Should be expired")
}

func TestIsExpired2(t *testing.T) {
	TimeNow := time.Now()
	time.Sleep(time.Second * 5)
	TokenCreatedTime := time.Now()
	assert.Equal(t, true, IsExpired(TimeNow, TokenCreatedTime, 3), "Should be expired")
}

func TestIsExpired3(t *testing.T) {
	TokenCreatedTime := time.Now()
	time.Sleep(time.Second * 3)
	TimeNow := time.Now()
	assert.Equal(t, false, IsExpired(TimeNow, TokenCreatedTime, 10), "Should not be expired")
}
