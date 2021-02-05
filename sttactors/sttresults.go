package sttactors

import (
	"fmt"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/asteriskclient"
	"github.com/va-voice-gateway/gateway"
	"github.com/va-voice-gateway/nlp"
	"github.com/va-voice-gateway/tts"
	"log"
	"sync"
)

// singleton based on pattern described here
// http://blog.ralch.com/tutorial/design-patterns/golang-singleton/
var (
	_instance *sttResultsActor
	_once sync.Once
)

type CommandErrorResult struct {
	ChannelId string
	Error error
}

type CommandPartialResult struct {
	ChannelId string
	Text string
}

type CommandFinalResult struct {
	ChannelId string
	Text string
}

// https://stackoverflow.com/questions/36870289/goroutines-and-channels-with-multiple-types
type sttResultsActor struct {
	CommandsChannel chan interface{} // truly terrible :( sad go story...
}

func STTResultsActor() *sttResultsActor {
	_once.Do(func() {
		_instance = &sttResultsActor{
			CommandsChannel: make(chan interface{}),
		}
	})
	return _instance
}

func (sttra *sttResultsActor) STTResultsActorProcessingLoop() {
	for command := range sttra.CommandsChannel {
		switch v := command.(type) {
			case CommandErrorResult:
				log.Printf("STTResultsActorProcessingLoop.CommandErrorResult  %v\n", v)
				go sttra.errorResult(v)
				break
			case CommandPartialResult:
				log.Printf("STTResultsActorProcessingLoop.CommandPartialResult  %v\n", v)
				go sttra.partialResult(v)
				break
			case CommandFinalResult:
				log.Printf("STTResultsActorProcessingLoop.CommandFinalResult  %v\n", v)
				go sttra.finalResult(v)
				break
			default:
				log.Printf("STTResultsActorProcessingLoop.Unknown type, ignoring  %v\n", v)
		}
	}
}

func (sttra *sttResultsActor) errorResult(cmdErrorResult CommandErrorResult) {
	// TBD
}

func (sttra *sttResultsActor) partialResult(cmdPartialResult CommandPartialResult) {
	// TBD
}

func (sttra *sttResultsActor) finalResult(cmdFinalResult CommandFinalResult) {
	gw := gateway.GatewayService()
	botId, lang := gw.GetBotIdLang(&cmdFinalResult.ChannelId)
	if botId != nil && lang != nil {
		Nlp_tts_play(&cmdFinalResult.ChannelId, botId, lang, nlp.NLPRequest{
			Text: &nlp.NLPRequestText {
				Text: cmdFinalResult.Text,
			},
			Event: nil,
		})
	}
}

// TBD: this will have to be probably public so that we can use from different places
// the challenge here are import cycles
func Nlp_tts_play(clientId *string, botId *string, language *string, nlpRequest nlp.NLPRequest) {
	appConfig := appconfig.AppConfig(nil)
	gw := gateway.GatewayService()

	// TBD: should CallNLP & InvokeTTS  be called as go routines ?

	nlpRes, err := gw.CallNLP(clientId, nlpRequest)
	if err != nil {
		log.Printf("Nlp_tts_play error(CallNLP) %s\n", err)
		return
	}

	ttsRes, err := tts.InvokeTTS(tts.TTSReq{
		BotId: *botId,
		Text: nlpRes.Text,
		Lang: *language,
	})
	if err != nil {
		log.Printf("Nlp_tts_play error(InvokeTTS) %s\n", err)
		return
	}

	log.Println("File to play " + ttsRes.FileName)

	aric := *asteriskclient.AriClient
	channelId := ari.NewKey(ari.ChannelKey, *clientId)
	playbackID := ""
	mediaURI := fmt.Sprintf("sound:%s%s", appConfig.Tts.TtsBaseUrlAsterisk,ttsRes.FileName)
	playbackHandle, err := aric.Channel().Play(channelId, playbackID, mediaURI)

	if err != nil {
		log.Printf("Nlp_tts_play error(Play) %s\n", err)
		return
	}

	log.Println("playback ", playbackHandle.ID())

}