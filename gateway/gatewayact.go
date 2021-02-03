package gateway

import (
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

type gatewayActor struct {
	CommandsChannel chan interface{}
	Gateway Gateway
}

func GatewayActor() *gatewayActor {
	once.Do(func() {
		instance = &gatewayActor {
			CommandsChannel: make(chan interface{}),
			Gateway: newGateway(),
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
			default:
				log.Printf("GatewayActorProcessingLoop.Unknown type, ignoring  %v\n", v)
		}
	}
}