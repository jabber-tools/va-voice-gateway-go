package appconfig

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	Tts Tts
	Nlp Nlp
	NlpVap NlpVap
	Asterisk Asterisk
	Core Core
	Log Log
	Temp Temp
}

type Tts struct {
	TtsBaseUrl string `toml:"tts_base_url"`
	TtsBaseUrlAsterisk string `toml:"tts_base_url_asterisk"`
	TtsApiBasicAuthToken string `toml:"tts_api_basic_auth_token"`
}

type Nlp struct {
	TargetNlp string `toml:"target_nlp"`
}

type NlpVap struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	VapBaseUrl string `toml:"vap_base_url"`
}

type Asterisk struct {
	AriUrl string `toml:"ari_url"`
	WSUrl string `toml:"ws_url"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	App string `toml:"app"`
}

type Core struct {
	Port int `toml:"port"`
	Host string `toml:"host"`
	ChannelSize int `toml:"channel_size"`
}

type Log struct {
	LogCfg string `toml:"log_cfg"`
}

type Temp struct {
	SttMsSubKey string `toml:"stt_ms_sub_key"`
	SttMsRegion string `toml:"stt_ms_region"`
}

func LoadAppConfig(cappCfgPath string) (*AppConfig, error) {
	fmt.Println("loading app config...")

	var config AppConfig
	if _, err := toml.DecodeFile(cappCfgPath, &config); err != nil {
		return nil, err
	}

	return &config, nil

}