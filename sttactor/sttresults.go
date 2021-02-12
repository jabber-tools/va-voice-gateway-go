package sttactor

import (
	"github.com/CyCoreSystems/ari/v5"
	"github.com/va-voice-gateway/asteriskclient"
	"github.com/va-voice-gateway/gateway"
	"github.com/va-voice-gateway/logger"
	"github.com/va-voice-gateway/nlp"
	"github.com/va-voice-gateway/utils"
	"github.com/sirupsen/logrus"
	"sync"
)

// singleton based on pattern described here
// http://blog.ralch.com/tutorial/design-patterns/golang-singleton/
var (
	_instance *sttResultsActor
	_once sync.Once
	log = logrus.New()
)

func init() {
	logger.InitLogger(log, "sttactor")
}

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
				log.Debugf("STTResultsActorProcessingLoop.CommandErrorResult  %v\n", v)
				go sttra.errorResult(v)
				break
			case CommandPartialResult:
				log.Debugf("STTResultsActorProcessingLoop.CommandPartialResult  %v\n", v)
				go sttra.partialResult(v)
				break
			case CommandFinalResult:
				log.Debugf("STTResultsActorProcessingLoop.CommandFinalResult  %v\n", v)
				go sttra.finalResult(v)
				break
			default:
				log.Debugf("STTResultsActorProcessingLoop.Unknown type, ignoring  %v\n", v)
		}
	}
}

func (sttra *sttResultsActor) errorResult(cmdErrorResult CommandErrorResult) {
	// TBD
}

func (sttra *sttResultsActor) partialResult(cmdPartialResult CommandPartialResult) {
	channelId := cmdPartialResult.ChannelId
	gw := gateway.GatewayService()
	if playbackId := gw.GetPlaybackId(&channelId); playbackId != nil {
		log.Debugf("Stopping playback %s for %s", playbackId, channelId)
		ariClient := *asteriskclient.AriClient
		ariClient.Playback().Stop(ari.NewKey(ari.PlaybackKey, *playbackId))
		gw.ResetPlaybackId(&channelId)
	}
}

func (sttra *sttResultsActor) finalResult(cmdFinalResult CommandFinalResult) {
	gw := gateway.GatewayService()
	botId, lang := gw.GetBotIdLang(&cmdFinalResult.ChannelId)
	if botId != nil && lang != nil {
		normalizedText := utils.RemoveNonAlphaNumericChars(utils.NormalizeAWB(cmdFinalResult.Text))
		log.Debugf("Normalized text %s\n", normalizedText)
		go asteriskclient.Nlp_tts_play(&cmdFinalResult.ChannelId, botId, lang, nlp.NLPRequest{
			Text: &nlp.NLPRequestText {
				Text: normalizedText,
			},
			Event: nil,
		})
	}
}

