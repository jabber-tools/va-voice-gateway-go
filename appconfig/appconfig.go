package appconfig

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/va-voice-gateway/logger"
	"sync"
)

var (
	instance *appConfig
	once sync.Once
	log = logrus.New()
)

func init() {
	logger.InitLogger(log)
}

type appConfig struct {
	Tts      Tts
	Nlp      Nlp
	NlpVap   NlpVap
	Asterisk Asterisk
	Core     Core
	Log      Log
	Temp     Temp
}

type Tts struct {
	TtsBaseUrl           string `toml:"tts_base_url"`
	TtsBaseUrlAsterisk   string `toml:"tts_base_url_asterisk"`
	TtsApiBasicAuthToken string `toml:"tts_api_basic_auth_token"`
}

type Nlp struct {
	TargetNlp string `toml:"target_nlp"`
}

type NlpVap struct {
	Username   string `toml:"username"`
	Password   string `toml:"password"`
	VapBaseUrl string `toml:"vap_base_url"`
}

type Asterisk struct {
	AriUrl   string `toml:"ari_url"`
	WSUrl    string `toml:"ws_url"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	App      string `toml:"app"`
}

type Core struct {
	Port        int    `toml:"port"`
	Host        string `toml:"host"`
	ChannelSize int    `toml:"channel_size"`
}

type Log struct {
	LogCfg string `toml:"log_cfg"`
}

type Temp struct {
	SttMsSubKey string `toml:"stt_ms_sub_key"`
	SttMsRegion string `toml:"stt_ms_region"`
}

func loadAppConfig(appCfgPath string) (*appConfig, error) {
	log.Info("loading app config...")

	var config appConfig
	if _, err := toml.DecodeFile(appCfgPath, &config); err != nil {
		return nil, err
	}

	return &config, nil

}

func AppConfig(appCfgPath *string) *appConfig {
	once.Do(func() {
		appConfig, err := loadAppConfig(*appCfgPath)
		if err != nil {
			log.Error("Error when loading app config")
			log.Fatal(err)
			return
		}
		instance = appConfig
	})
	return instance
}
