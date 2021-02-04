package gateway

import "github.com/va-voice-gateway/nlp"

// TBD: rust version contains other attributes
// add them on the fly as needed
type Client struct {
	ClientId     string
	PlaybackId   *string
	DoSTT        bool
	Terminate    bool
	Dtmf         []string
	BotId        string
	Lang         string
	InviteParams map[string]string
	NLP 	     nlp.VAP
}

func NewClient(clientId string, botId string, lang string, inviteParams map[string]string) Client {
	return Client{
		ClientId:     clientId,
		PlaybackId:   nil,
		DoSTT:        false, // TBD: conditionally true for audio buffering feature
		Terminate:    false,
		Dtmf:         make([]string, 0),
		BotId:        botId,
		Lang:         lang,
		InviteParams: inviteParams,
	}
}
