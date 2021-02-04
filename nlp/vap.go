package nlp

import (
	"encoding/json"
	"fmt"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/gateway/config"
	"github.com/va-voice-gateway/utils"
	"log"
	"github.com/google/uuid"
)

type VAP struct {
	NewConv bool
	VapAccessToken string
	ClientId string
	BotId string
	Lang string
	InviteParams string
	ConvId string
}

type NLPRequestText struct {
	Text string
}

type NLPRequestEvent struct {
	Name string
	Params interface{}
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
			ConvId: uuid.NewString(),
		}, nil
	} else {
		return nil, fmt.Errorf("NewVAP: bot config not found for %s", BotId)
	}
}

func (v *VAP) InvokeNLP(request *NLPRequest) (*NLPResponse, error) {
	var payload string
	token := utils.GetVapAPIToken()
	appConfig := appconfig.AppConfig()

	if request.Text != nil /* text nlp request */ {
		if v.NewConv == true /* new conv */ {
			v.NewConv = false
			payload = fmt.Sprintf(`{
				"headers": {
					"at": "%s"
				},
				"body": {
					"text": "%s",
					"convId": "%s",
				},
				"vaContext": {
					"lang": "%s",
					"voicegw": {
						inviteParams: "%s"
					}			
				}
			}`, v.VapAccessToken, request.Text.Text, v.ConvId, v.Lang, v.InviteParams)

		} else /* ongoing conv */ {
			payload = fmt.Sprintf(`{
				"headers": {
					"at": "%s"
				},
				"body": {
					"text": "%s",
					"convId": "%s",
				}
			}`, v.VapAccessToken, request.Text.Text, v.ConvId)
		}
	} else /* event nlp request */ {
		if v.NewConv == true /* new conv */ {
			v.NewConv = false
			payload = fmt.Sprintf(`{
				"headers": {
					"at": "%s"
				},
				"body": {
					"event": {
 						"name": "%s"
					},
					"convId": "%s",
				},
				"vaContext": {
					"lang": "%s",
					"voicegw": {
						inviteParams: "%s"
					}			
				}
			}`, v.VapAccessToken, request.Event.Name, v.ConvId, v.Lang, v.InviteParams)

		} else /* ongoing conv */ {
			payload = fmt.Sprintf(`{
				"headers": {
					"at": "%s"
				},
				"body": {
					"event": {
 						"name": "%s"
					},
					"convId": "%s",
				}
			}`, v.VapAccessToken, request.Event.Name, v.ConvId)
		}
	}

	log.Printf("payload %s\n", payload)
	log.Printf("token %s\n", token)

	url := fmt.Sprintf("%s%s", appConfig.NlpVap.VapBaseUrl, "/vapapi/authentication/v1")
	log.Printf("url %s\n", url)

	return nil, nil
}

