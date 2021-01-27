package gateway

type BotConfig struct {

}

type Gateway struct {
	Clients map[string]Client
	BotConfigs []BotConfig
}

func NewGateway(botConfigs []BotConfig) Gateway {
	return Gateway{
		Clients: make(map[string]Client),
		BotConfigs: botConfigs,
	}
}

func (g *Gateway) AddClient(client Client) {
	g.Clients[client.ClientId] = client
}

func (g *Gateway) RemoveClient(client Client) {
	delete(g.Clients, client.ClientId)
}