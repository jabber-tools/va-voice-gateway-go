package asteriskclient

import (
	"fmt"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/gateway"
	"github.com/va-voice-gateway/logger"
	"github.com/va-voice-gateway/nlp"
	"github.com/va-voice-gateway/tts"
	"github.com/sirupsen/logrus"
	"regexp"
)

var (
	AriClient *ari.Client
	ReNewLineChar = regexp.MustCompile(`\n`)
	log = logrus.New()
)

func init() {
	logger.InitLogger(log, "asteriskclient")
}

// helper composite function to perform nlp, tts and asterisk play
// put into asteriskclient package to prevent import cycle though ideally it should go to asterisk package
func Nlp_tts_play(clientId *string, botId *string, language *string, nlpRequest nlp.NLPRequest) {
	appConfig := appconfig.AppConfig(nil)
	gw := gateway.GatewayService()

	// TBD: should CallNLP & InvokeTTS  be called as go routines ?

	nlpRes, err := gw.CallNLP(clientId, nlpRequest)
	if err != nil {
		log.Error("Nlp_tts_play error(CallNLP) %s\n", err)
		return
	}

	escapedText := ReNewLineChar.ReplaceAllString(nlpRes.Text, "\\n")
	log.Debug("Invoking TTS %s\n", escapedText)

	ttsRes, err := tts.InvokeTTS(tts.TTSReq{
		BotId: *botId,
		Text: escapedText,
		Lang: *language,
	})
	if err != nil {
		log.Error("Nlp_tts_play error(InvokeTTS) %s\n", err)
		return
	}

	log.Info("File to play " + ttsRes.FileName)

	aric := *AriClient
	channelId := ari.NewKey(ari.ChannelKey, *clientId)
	playbackID := ""
	mediaURI := fmt.Sprintf("sound:%s%s", appConfig.Tts.TtsBaseUrlAsterisk,ttsRes.FileName)
	playbackHandle, err := aric.Channel().Play(channelId, playbackID, mediaURI)

	if err != nil {
		log.Error("Nlp_tts_play error(Play) %s\n", err)
		return
	}

	playbackId := playbackHandle.ID()
	gw.SetPlaybackId(clientId, &playbackId)
	log.Debug("playback ",playbackId)
}
