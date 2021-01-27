package gateway

// TBD: rust version contains other attributes
// we probably do not need them in Go solution
type Client struct {
	ClientId string
	PlaybackId string
	DoSTT bool
	Dtmf []string
	BotId string
	Lang string
	InviteParams map[string]string
}

func NewClient(clientId string, botId string, lang string, inviteParams map[string]string) Client {
	return Client {
		ClientId: clientId,
		PlaybackId: nil,
		DoSTT: false, // TBD: conditionally true for audio buffering feature
		Dtmf: make([]string, 0),
		BotId: botId,
		Lang: lang,
		InviteParams: inviteParams,
	}
}