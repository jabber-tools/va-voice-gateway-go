package sttactors

import (
	"github.com/va-voice-gateway/gateway"
	"github.com/va-voice-gateway/nlp"
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
func Nlp_tts_play(clientId *string, botId *string, language * string, nlpRequest nlp.NLPRequest) {
	gw := gateway.GatewayService()
	gw.CallNLP(clientId, nlpRequest)
}