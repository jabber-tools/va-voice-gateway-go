package actors

import (
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
				break
			case CommandPartialResult:
				log.Printf("STTResultsActorProcessingLoop.CommandPartialResult  %v\n", v)
				break
			case CommandFinalResult:
				log.Printf("STTResultsActorProcessingLoop.CommandFinalResult  %v\n", v)
				break
			default:
				log.Printf("STTResultsActorProcessingLoop.Unknown type, ignoring  %v\n", v)
		}
	}
}