package gateway

import (
	"github.com/va-voice-gateway/nlp"
	"sync"
)

var (
	_instance *gatewaySvc
	_once sync.Once
)

type gatewaySvc struct {
	// just convenience so that we dont
	// have to call GatewayActor() every time
	GWActor *gatewayActor
}

func GatewayService() *gatewaySvc {
	_once.Do(func() {
		_instance = &gatewaySvc {
			GWActor: GatewayActor(),
		}
	})
	return _instance
}

func (g *gatewaySvc) AddClient(client Client) {
	g.GWActor.CommandsChannel <- CommandAddClient{
		Client: client,
	}
}

func (g *gatewaySvc) RemoveClient(clientId string) {
	g.GWActor.CommandsChannel <- CommandRemoveClient{
		ClientId: clientId,
	}
}

func (g *gatewaySvc) SetPlaybackId(clientId *string, playbackId *string) {
	g.GWActor.CommandsChannel <- CommandSetPlaybackId{
		ClientId: *clientId,
		PlaybackId: *playbackId,
	}
}

func (g *gatewaySvc) GetPlaybackId(clientId *string) *string {
	c := make(chan *string)
	g.GWActor.CommandsChannel <- CommandGetPlaybackId{
		ClientId: *clientId,
		Responder: c,
	}
	return <- c
}

func (g *gatewaySvc) ResetPlaybackId(clientId *string) {
	g.GWActor.CommandsChannel <- CommandResetPlaybackId{
		ClientId: *clientId,
	}
}

func (g *gatewaySvc) SetTerminating(clientId *string) {
	g.GWActor.CommandsChannel <- CommandSetIsTerminating{
		ClientId: *clientId,
	}
}

func (g *gatewaySvc) GetTerminating(clientId *string) bool {
	c := make(chan bool)
	g.GWActor.CommandsChannel <- CommandGetIsTerminating{
		ClientId: *clientId,
		Responder: c,
	}
	return <- c
}

func (g *gatewaySvc) SetDoSTT(clientId *string, doSTT bool) {
	g.GWActor.CommandsChannel <- CommandSetDoSTT{
		ClientId: *clientId,
		DoSTT: doSTT,
	}
}

func (g *gatewaySvc) GetDoSTT(clientId *string) bool {
	c := make(chan bool)
	g.GWActor.CommandsChannel <- CommandGetDoSTT{
		ClientId: *clientId,
		Responder: c,
	}
	return <- c
}

func (g *gatewaySvc) GetBotIdLang(clientId *string) (*string, *string) {
	c := make(chan BotIdLang)
	g.GWActor.CommandsChannel <- CommandGetBotIdLang{
		ClientId: *clientId,
		Responder: c,
	}
	botidlang := <- c
	return botidlang.BotId, botidlang.Lang
}

func (g *gatewaySvc) AddDtmf(clientId *string, dtmf string) {
	g.GWActor.CommandsChannel <- CommandAddDtmf{
		ClientId: *clientId,
		Dtmf: dtmf,
	}
}

func (g *gatewaySvc) GetDtmf(clientId *string) string {
	c := make(chan string)
	g.GWActor.CommandsChannel <- CommandGetDtmf{
		ClientId: *clientId,
		Responder: c,
	}
	return <- c
}

func (g *gatewaySvc) ResetDtmf(clientId *string) {
	g.GWActor.CommandsChannel <- CommandResetDtmf{
		ClientId: *clientId,
	}
}

func (g *gatewaySvc) CallNLP(clientId *string, nlpRequest nlp.NLPRequest) (*nlp.NLPResponse, error) {
	c := make(chan nlp.NLPResponseResult)
	g.GWActor.CommandsChannel <- CommandCallNLP{
		ClientId: *clientId,
		Request: nlpRequest,
		Responder: c,
	}
	result := <- c
	if result.NLPResponse != nil {
		return result.NLPResponse, nil
	} else {
		return nil, result.Error
	}
}