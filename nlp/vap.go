package nlp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/gateway/config"
	"github.com/va-voice-gateway/utils"
	"io/ioutil"
	"log"
	"net/http"
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

type NLPResponseResult struct {
	NLPResponse *NLPResponse
	Error error
}

type DiagnosticsInfo struct {
	EndConversation bool `json:"end_conversation"`
}

type QueryResult struct {
	DiagnosticsInfo *DiagnosticsInfo `json:"diagnosticInfo"`
}

type DialogFlowResponse struct {
	QueryResult *QueryResult `json:"queryResult"`
}

type VAPCanonicalResponse struct {
	DfResponse *DialogFlowResponse `json:"dfResponse"`
}

type VAPResponse struct {
	CanonicalResponse VAPCanonicalResponse `json:"canonicalResponse"`
	VoiceGWResponse string `json:"voiceGwResponse"`
}

func NewVAP(ClientId string, BotId string, Lang string, InviteParams map[string]string) (*VAP, error) {

	jsonStringInviteParams, err := json.Marshal(InviteParams)
	if err != nil {
		return nil, err
	}

	botConfigs := config.BotConfigs(nil)
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
	appConfig := appconfig.AppConfig(nil)

	if request.Text != nil /* text nlp request */ {
		if v.NewConv == true /* new conv */ {
			v.NewConv = false
			payload = fmt.Sprintf(`{
				"headers": {
					"at": "%s"
				},
				"body": {
					"text": "%s",
					"convId": "%s"
				},
				"vaContext": {
					"lang": "%s",
					"voicegw": {
						"inviteParams": %s
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
					"convId": "%s"
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
					"convId": "%s"
				},
				"vaContext": {
					"lang": "%s",
					"voicegw": {
						"inviteParams": %s
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
					"convId": "%s"
				}
			}`, v.VapAccessToken, request.Event.Name, v.ConvId)
		}
	}

	// log.Printf("payload %s\n", payload)
	// log.Printf("token %s\n", *token)

	url := fmt.Sprintf("%s%s", appConfig.NlpVap.VapBaseUrl, "/vapapi/channels/voicegw/v1")


	client := &http.Client{}
	vapToken := utils.GetVapAPIToken()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", *vapToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Printf("InvokeNLP: error when calling  /vapapi/channels/voicegw/v1: %v\n", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("InvokeNLP: error when reading http response: %v\n", err)
		return nil, err
	}

	log.Printf("InvokeNLP: raw response: %v\n", string(body))

	vapResponse := &VAPResponse{}
	err = json.Unmarshal(body, vapResponse)
	if err != nil {
		log.Printf("InvokeNLP: error when parsing json: %v\n", err)
		return nil, err
	}

	vapResponseStr,_ := utils.StructToJsonString(vapResponse)
	log.Printf("InvokeNLP: response parsed: %s\n", *vapResponseStr)

	var IsEOC bool
	if vapResponse.CanonicalResponse.DfResponse != nil &&
	   vapResponse.CanonicalResponse.DfResponse.QueryResult != nil &&
	   vapResponse.CanonicalResponse.DfResponse.QueryResult.DiagnosticsInfo != nil &&
		vapResponse.CanonicalResponse.DfResponse.QueryResult.DiagnosticsInfo.EndConversation == true {
		IsEOC = true
	} else {
		IsEOC = false
	}

	return &NLPResponse {
		Text: vapResponse.VoiceGWResponse,
		IsEOC: IsEOC,
	}, nil
}

