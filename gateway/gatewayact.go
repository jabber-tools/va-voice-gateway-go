package gateway

// in order to prevent cyclical dependencies gateway actor
// is defined in gateway package and
// not in actors package as remaining actors

import (
	"github.com/va-voice-gateway/nlp"
	"log"
	"sync"
)

var (
	instance *gatewayActor
	once sync.Once
)

type CommandAddClient struct {
	Client Client
}

type CommandRemoveClient struct {
	ClientId string
}

type CommandSetPlaybackId struct {
	ClientId string
	PlaybackId string
}

type CommandGetPlaybackId struct {
	ClientId string
	Responder chan *string
}

type CommandResetPlaybackId struct {
	ClientId string
}

type CommandGetIsTerminating struct {
	ClientId string
	Responder chan bool
}

type CommandSetIsTerminating struct {
	ClientId string
}

type CommandGetDoSTT struct {
	ClientId string
	Responder chan bool
}

type CommandSetDoSTT struct {
	ClientId string
	DoSTT bool
}

// helper struct to bypass fact
// go does not support tuples
type BotIdLang struct {
	BotId *string
	Lang *string
}

type CommandGetBotIdLang struct {
	ClientId string
	Responder chan BotIdLang
}

type CommandAddDtmf struct {
	ClientId string
	Dtmf string
}

type CommandGetDtmf struct {
	ClientId string
	Responder chan string
}

type CommandResetDtmf struct {
	ClientId string
}

type CommandCallNLP struct {
	ClientId string
	Request nlp.NLPRequest
	Responder chan nlp.NLPResponseResult
}

type gatewayActor struct {
	CommandsChannel chan interface{}
	Gateway         Gateway
}

func GatewayActor() *gatewayActor {
	once.Do(func() {
		instance = &gatewayActor{
			CommandsChannel: make(chan interface{}),
			Gateway:         newGateway(),
		}
	})
	return instance
}

func (gwa *gatewayActor) GatewayActorProcessingLoop() {
	for command := range gwa.CommandsChannel {
		switch v := command.(type) {
			case CommandAddClient:
				gwa.Gateway.AddClient(v.Client)
				break
			case CommandRemoveClient:
				gwa.Gateway.RemoveClient(v.ClientId)
				break
			case CommandSetPlaybackId:
				gwa.Gateway.ClientSetPlaybackId(&v.ClientId, &v.PlaybackId)
				break
			case CommandGetPlaybackId:
				playbackId := gwa.Gateway.ClientGetPlaybackId(&v.ClientId)
				v.Responder <- playbackId
				break
			case CommandResetPlaybackId:
				gwa.Gateway.ClientResetPlaybackId(&v.ClientId)
				break
			case CommandGetIsTerminating:
				isTerminating := gwa.Gateway.ClientGetTerminating(&v.ClientId)
				v.Responder <- isTerminating
				break
			case CommandSetIsTerminating:
				gwa.Gateway.ClientSetTerminating(&v.ClientId)
				break
			case CommandGetDoSTT:
				isTerminating := gwa.Gateway.ClientGetDoSTT(&v.ClientId)
				v.Responder <- isTerminating
				break
			case CommandSetDoSTT:
				gwa.Gateway.ClientSetDoSTT(&v.ClientId, v.DoSTT)
				break
			case CommandGetBotIdLang:
				botId, lang := gwa.Gateway.ClientGetBotIdLang(&v.ClientId)
				v.Responder <- BotIdLang{
					BotId: botId,
					Lang: lang,
				}
				break
			case CommandAddDtmf:
				gwa.Gateway.ClientAddDtmf(&v.ClientId, v.Dtmf)
				break
			case CommandGetDtmf:
				dtmf := gwa.Gateway.ClientGetDtmf(&v.ClientId)
				v.Responder <- *dtmf
				break
			case CommandResetDtmf:
				gwa.Gateway.ClientResetDtmf(&v.ClientId)
				break

			case CommandCallNLP:
				go func() {
					clientNLP := gwa.Gateway.ClientGetNLP(&v.ClientId)
					if clientNLP != nil {
						nlpResponse, err := clientNLP.InvokeNLP(&v.Request)
						if err != nil {
							v.Responder <- nlp.NLPResponseResult {
								NLPResponse: nil,
								Error: err,
							}
						} else {
							v.Responder <- nlp.NLPResponseResult {
								NLPResponse: nlpResponse,
								Error: nil,
							}
						}
					}
				}()
				break


			default:
				log.Printf("GatewayActorProcessingLoop.Unknown type, ignoring  %v\n", v)
		}
	}
}