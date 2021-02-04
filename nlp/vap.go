package nlp

import (
	"encoding/json"
	"fmt"
	"github.com/va-voice-gateway/gateway/config"
)

type VAP struct {
	NewConv bool
	VapAccessToken string
	ClientId string
	BotId string
	Lang string
	InviteParams string
}

type NLPRequestText struct {
	Text string
}

type NLPRequestEvent struct {
	EventName string
	EventParams interface{}
}

// since go does not support enums with embedded structs
// we model it like two pointers (since go does not support optional types as well:) )
// one of the pointers will be set the other nil. they cannot be set both!
type NLPRequest struct {
	Text  *NLPRequestText
	Event *NLPRequestEvent
}

type NLPResponse struct {
	Text string
	IsEOC bool
}

func NewVAP(ClientId string, BotId string, Lang string, InviteParams map[string]string) (*VAP, error) {

	jsonStringInviteParams, err := json.Marshal(InviteParams)
	if err != nil {
		return nil, err
	}

	botConfigs := config.BotConfigs()
	botConfig := botConfigs.GetBotConfig(&BotId)
	if botConfig != nil {
		return &VAP {
			NewConv: true,
			VapAccessToken: botConfig.Channels.Webchat.AccessToken,
			ClientId:  ClientId,
			BotId: BotId,
			Lang:  Lang,
			InviteParams: string(jsonStringInviteParams),
		}, nil
	} else {
		return nil, fmt.Errorf("NewVAP: bot config not found for %s", BotId)
	}
}

func (v *VAP) InvokeNLP(request *NLPRequest) (*NLPResponse, error) {
	// token := utils.GetVapAPIToken()

	return nil, nil
}

