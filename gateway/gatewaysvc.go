package gateway

import (
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
	g.GWActor.CommandsChannel <- CommandAddClient {
		Client: client,
	}
}

func (g *gatewaySvc) RemoveClient(clientId string) {
	g.GWActor.CommandsChannel <- CommandRemoveClient {
		ClientId: clientId,
	}
}

func (g *gatewaySvc) SetPlaybackId(clientId *string, playbackId *string) {
	g.GWActor.CommandsChannel <- CommandSetPlaybackId {
		ClientId: *clientId,
		PlaybackId: *playbackId,
	}
}

func (g *gatewaySvc) GetPlaybackId(clientId *string) *string {
	c := make(chan *string)
	g.GWActor.CommandsChannel <- CommandGetPlaybackId {
		ClientId: *clientId,
		Responder: c,
	}
	return <- c
}

