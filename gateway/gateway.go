package gateway

import "strings"

type BotConfig struct {

}

type Gateway struct {
	Clients map[*string]Client
	BotConfigs []BotConfig
}

func NewGateway(botConfigs []BotConfig) Gateway {
	return Gateway{
		Clients: make(map[*string]Client),
		BotConfigs: botConfigs,
	}
}

func (g *Gateway) AddClient(client Client) {
	g.Clients[&client.ClientId] = client
}

func (g *Gateway) RemoveClient(client Client) {
	delete(g.Clients, &client.ClientId)
}

func (g *Gateway) ClientSetPlaybackId(clientId *string, playbackId *string) {
	if client, ok := g.Clients[clientId]; ok {
		client.PlaybackId = playbackId
	}
}

func (g *Gateway) ClientGetPlaybackId(clientId *string) *string {
	if client, ok := g.Clients[clientId]; ok {
		return client.PlaybackId
	} else {
		return nil
	}
}

func (g *Gateway) ClientResetPlaybackId(clientId *string) {
	if client, ok := g.Clients[clientId]; ok {
		client.PlaybackId = nil
	}
}

func (g *Gateway) ClientSetTerminating(clientId *string) {
	if client, ok := g.Clients[clientId]; ok {
		client.Terminate = true
	}
}

func (g *Gateway) ClientGetTerminating(clientId *string) bool {
	if client, ok := g.Clients[clientId]; ok {
		return client.Terminate
	} else {
		return true
	}
}

func (g *Gateway) ClientSetDoSTT(clientId *string, doSTT bool) {
	if client, ok := g.Clients[clientId]; ok {
		client.DoSTT = doSTT
	}
}

func (g *Gateway) ClientGetDoSTT(clientId *string) bool {
	if client, ok := g.Clients[clientId]; ok {
		return client.DoSTT
	} else {
		return false
	}
}

func (g *Gateway) ClientGetBotIdLang(clientId *string) (*string, *string) {
	if client, ok := g.Clients[clientId]; ok {
		return &client.BotId, &client.Lang
	} else {
		return nil, nil
	}
}

func (g *Gateway) ClientAddDtmf(clientId *string, val string) {
	if client, ok := g.Clients[clientId]; ok {
		client.Dtmf = append(client.Dtmf, val)
	}
}

func (g *Gateway) ClientResetDtmf(clientId *string) {
	if client, ok := g.Clients[clientId]; ok {
		client.Dtmf = make([]string, 0)
	}
}

func (g *Gateway) ClientGetDtmf(clientId *string) *string {
	if client, ok := g.Clients[clientId]; ok {
		dtmfs := strings.Join(client.Dtmf, "")
		return &dtmfs
	}
	return nil
}

