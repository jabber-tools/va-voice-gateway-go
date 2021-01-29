package gateway

import (
	"github.com/va-voice-gateway/stt/google"
	"strings"
)

type BotConfig struct {
	Name string `json:"name"`
	BotId string `json:"botId"`
	Channels BotConfigChannels `json:"channels"`
	Brain BotConfigBrain `json:"brain"`
}

type BotConfigChannels struct {
	VoiceGW BotConfigChannelsVoiceGW `json:"voicegw"`
}

type BotConfigChannelsVoiceGW struct {
	TTSContentUrl string `json:"ttsContentUrl"`
	Mapping []BotConfigChannelsVoiceGWMapping `json:"mapping"`
	Providers BotConfigChannelsVoiceGWProviders `json:"providers"`
}


type BotConfigChannelsVoiceGWMapping struct {
	Lang string `json:"lang"`
	TTSProvider string `json:"ttsProvider"`
	STTProvider string `json:"sttProvider"`
}

type BotConfigChannelsVoiceGWProviders struct {
	Google GoogleProvider `json:"google"`
	Microsoft MicrosoftProvider `json:"microsoft"`
}

type GoogleProvider struct {
	TTSApiUrl string `json:"ttsApiUrl"`
	Credentials GDFCredentials `json:"credentials"`
	STT []LanguageConfig `json:"stt"`
}

type LanguageConfig struct {
	Lang string `json:"lang"`
	Cfg GoogleConfigItemCfg `json:"cfg"`
}

type GoogleConfigItemCfg struct {
	SpeechRecognitionConfig google.RecognitionConfig `json:"speechRecognitionConfig"`
}

type MicrosoftProvider struct {
	TTSApiUrl string `json:"ttsApiUrl"`
	Credentials MicrosoftProviderCredentials `json:"credentials"`
}

type MicrosoftProviderCredentials struct {
	SubscriptionKey string `json:"subscriptionKey"`
	Region string `json:"region"`
}


type BotConfigBrain struct {
	Provider string `json:"provider"`
	Dialogflow BotConfigBrainDialogflow `json:"dialogflow"`
}

type BotConfigBrainDialogflow struct {
	Credentials GDFCredentials `json:"credentials"`
}

type GDFCredentials struct {
	Type string `json:"type"`
	ProjectId string `json:"project_id"`
	PrivateKeyId string `json:"private_key_id"`
	PrivateKey string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	ClientId string `json:"client_id"`
	AuthUri string `json:"auth_uri"`
	TokenUri string  `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl string  `json:"client_x509_cert_url"`
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

func (g *Gateway) GetBotConfig(botId *string) *BotConfig {
	for _, bc := range g.BotConfigs {
		if bc.BotId == *botId {
			return &bc
		}
	}
	return nil
}

// TBD: right now working with google STT provider only should be renamed
// same flaw has Rust version
func (g *Gateway) GetSTTBotConfig(botId *string, lang *string) *google.RecognitionConfig {
	botConfig := g.GetBotConfig(botId)
	if botConfig != nil {
		for _, langConfig := range botConfig.Channels.VoiceGW.Providers.Google.STT {
			if langConfig.Lang == *lang {
				return &langConfig.Cfg.SpeechRecognitionConfig
			}
		}
	}
	return nil
}



