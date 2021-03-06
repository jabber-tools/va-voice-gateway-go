package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/logger"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.New()

func init() {
	logger.InitLogger(log, "tts")
}

type TTSReq struct {
	BotId string `json:"bot_id"`
	Text string `json:"text"`
	Lang string `json:"lang"`
}

type TTSRes struct {
	FileName string `json:"fileName"`
}

func InvokeTTS(ttsReq TTSReq) (*TTSRes, error) {

	appConfig := appconfig.AppConfig(nil)

	payload := fmt.Sprintf(`{
					"bot_id": "%s",
					"text": "%s",
					"lang": "%s"
	}`, ttsReq.BotId, ttsReq.Text, ttsReq.Lang)

	log.Printf("InvokeTTS: request payload: %v\n", string(payload))

	url := fmt.Sprintf("%s%s", appConfig.Tts.TtsBaseUrl, "/tts/google/v1")

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s%s", "Basic ", appConfig.Tts.TtsApiBasicAuthToken))
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Printf("InvokeTTS: error when calling  /tts/google/v1: %v\n", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("InvokeTTS: error when calling  /tts/google/v1(wrong status code): %d", resp.StatusCode)
		return nil, fmt.Errorf("InvokeTTS: error when calling  /tts/google/v1(wrong status code): %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("InvokeTTS: error when reading http response: %v\n", err)
		return nil, err
	}

	log.Debugf("InvokeTTS: raw response: %v\n", string(body))

	ttsRes := &TTSRes{}
	err = json.Unmarshal(body, ttsRes)
	if err != nil {
		log.Errorf("InvokeTTS: error when parsing json: %v\n", err)
		return nil, err
	}

	return ttsRes, nil
}