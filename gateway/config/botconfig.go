package config

// Structs below are used to parse REST based Google STT configuration which are used by VAP
// These structs then needs to be transformed into respective GRPC based structs used by underlying Google APIs
type RecognitionConfig struct {
	Encoding                            int32                     `json:"encoding,omitempty"`
	SampleRateHertz                     int32                     `json:"sampleRateHertz,omitempty"`
	AudioChannelCount                   int32                     `json:"audioChannelCount,omitempty"`
	EnableSeparateRecognitionPerChannel bool                      `json:"enableSeparateRecognitionPerChannel,omitempty"`
	LanguageCode                        string                    `json:"languageCode,omitempty"`
	MaxAlternatives                     int32                     `json:"maxAlternatives,omitempty"`
	ProfanityFilter                     bool                      `json:"profanityFilter,omitempty"`
	SpeechContexts                      []*SpeechContext          `json:"speechContexts,omitempty"`
	EnableWordTimeOffsets               bool                      `json:"enableWordTimeOffsets,omitempty"`
	EnableAutomaticPunctuation          bool                      `json:"enableAutomaticPunctuation,omitempty"`
	DiarizationConfig                   *SpeakerDiarizationConfig `json:"diarizationConfig,omitempty"`
	Metadata                            *RecognitionMetadata      `json:"metadata,omitempty"`
	Model                               string                    `json:"model,omitempty"`
	UseEnhanced                         bool                      `json:"useEnhanced,omitempty"`
}

type RecognitionMetadata struct {
	InteractionType          int32  `json:"interactionType,omitempty"`
	IndustryNaicsCodeOfAudio uint32 `json:"industryNaicsCodeOfAudio,omitempty"`
	MicrophoneDistance       int32  `json:"microphoneDistance,omitempty"`
	OriginalMediaType        int32  `json:"originalMediaType,omitempty"`
	RecordingDeviceType      int32  `json:"recordingDeviceType,omitempty"`
	RecordingDeviceName      string `json:"recordingDeviceName,omitempty"`
	OriginalMimeType         string `json:"originalMimeType,omitempty"`
	AudioTopic               string `json:"audioTopic,omitempty"`
}

type SpeakerDiarizationConfig struct {
	EnableSpeakerDiarization bool  `json:"enableSpeakerDiarization,omitempty"`
	MinSpeakerCount          int32 `json:"minSpeakerCount,omitempty"`
	MaxSpeakerCount          int32 `json:"maxSpeakerCount,omitempty"`
	SpeakerTag               int32 `json:"speakerTag,omitempty"`
}

// TBD: are we missing boost attribute here (same in Rust version)?
type SpeechContext struct {
	Phrases []string `json:"phrases,omitempty"`
}

type BotConfig struct {
	Name     string            `json:"name"`
	BotId    string            `json:"botId"`
	Channels BotConfigChannels `json:"channels"`
	Brain    BotConfigBrain    `json:"brain"`
}

type BotConfigChannels struct {
	VoiceGW BotConfigChannelsVoiceGW `json:"voicegw"`
}

type BotConfigChannelsVoiceGW struct {
	TTSContentUrl string                            `json:"ttsContentUrl"`
	Mapping       []BotConfigChannelsVoiceGWMapping `json:"mapping"`
	Providers     BotConfigChannelsVoiceGWProviders `json:"providers"`
}

type BotConfigChannelsVoiceGWMapping struct {
	Lang        string `json:"lang"`
	TTSProvider string `json:"ttsProvider"`
	STTProvider string `json:"sttProvider"`
}

type BotConfigChannelsVoiceGWProviders struct {
	Google    GoogleProvider    `json:"google"`
	Microsoft MicrosoftProvider `json:"microsoft"`
}

type GoogleProvider struct {
	TTSApiUrl   string           `json:"ttsApiUrl"`
	Credentials GDFCredentials   `json:"credentials"`
	STT         []LanguageConfig `json:"stt"`
}

type LanguageConfig struct {
	Lang string              `json:"lang"`
	Cfg  GoogleConfigItemCfg `json:"cfg"`
}

type GoogleConfigItemCfg struct {
	SpeechRecognitionConfig RecognitionConfig `json:"speechRecognitionConfig"`
}

type MicrosoftProvider struct {
	TTSApiUrl   string                       `json:"ttsApiUrl"`
	Credentials MicrosoftProviderCredentials `json:"credentials"`
}

type MicrosoftProviderCredentials struct {
	SubscriptionKey string `json:"subscriptionKey"`
	Region          string `json:"region"`
}

type BotConfigBrain struct {
	Provider   string                   `json:"provider"`
	Dialogflow BotConfigBrainDialogflow `json:"dialogflow"`
}

type BotConfigBrainDialogflow struct {
	Credentials GDFCredentials `json:"credentials"`
}

type GDFCredentials struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}
type BotConfigs struct {
	BotConfigs []BotConfig
}

func NewBotConfigs(botConfigs []BotConfig) BotConfigs {
	return BotConfigs{
		BotConfigs: botConfigs,
	}
}

func (bcs *BotConfigs) GetBotConfig(botId *string) *BotConfig {
	for _, bc := range bcs.BotConfigs {
		if bc.BotId == *botId {
			return &bc
		}
	}
	return nil
}

// TBD: right now working with google STT provider only should be renamed
// same flaw has Rust version
func (bcs *BotConfigs) GetSTTBotConfig(botId *string, lang *string) *RecognitionConfig {
	botConfig := bcs.GetBotConfig(botId)
	if botConfig != nil {
		for _, langConfig := range botConfig.Channels.VoiceGW.Providers.Google.STT {
			if langConfig.Lang == *lang {
				return &langConfig.Cfg.SpeechRecognitionConfig
			}
		}
	}
	return nil
}

func (bcs *BotConfigs) GetSTTGoogleCred(botId *string) *GDFCredentials {
	if botConfig := bcs.GetBotConfig(botId); botConfig != nil {
		return &botConfig.Channels.VoiceGW.Providers.Google.Credentials
	}
	return nil
}

func (bcs *BotConfigs) GetNlpDialogflowCred(botId *string) *GDFCredentials {
	if botConfig := bcs.GetBotConfig(botId); botConfig != nil {
		return &botConfig.Brain.Dialogflow.Credentials
	}
	return nil
}