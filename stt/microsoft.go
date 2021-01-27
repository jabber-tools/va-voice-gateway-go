// +build linux

package stt

import (
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"github.com/va-voice-gateway/appconfig"
)

func PerformMicrosoftSTT(appConfig *appconfig.AppConfig, audioStream chan []byte) {
	audioFormat, err := audio.GetWaveFormatPCM(8000, 16, 1)
	if err != nil {
		fmt.Println("GetWaveFormatPCM error: ", err)
		return
	}
	defer audioFormat.Close()

	stream, err := audio.CreatePushAudioInputStreamFromFormat(audioFormat)
	if err != nil {
		fmt.Println("CreatePushAudioInputStreamFromFormat error: ", err)
		return
	}
	defer stream.Close()

	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
	if err != nil {
		fmt.Println("NewAudioConfigFromStreamInput error: ", err)
		return
	}
	defer audioConfig.Close()

	speechConfig, err := speech.NewSpeechConfigFromSubscription(appConfig.Temp.SttMsSubKey, appConfig.Temp.SttMsRegion)
	if err != nil {
		fmt.Println("NewAudioConfigFromStreamInput error: ", err)
		return
	}
	defer speechConfig.Close()

	speechRecognizer, err := speech.NewSpeechRecognizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		fmt.Println("NewSpeechRecognizerFromConfig error: ", err)
		return
	}
	defer speechRecognizer.Close()

	sessionStartedHandler := func(event speech.SessionEventArgs) {
		fmt.Println("sessionStartedHandler", event)
		defer event.Close()
	}

	sessionStoppedHandler := func(event speech.SessionEventArgs) {
		fmt.Println("sessionStoppedHandler", event)
		defer event.Close()
	}

	speechStartDetectedHandler := func(event speech.RecognitionEventArgs) {
		fmt.Println("speechStartDetectedHandler", event)
		defer event.Close()
	}

	speechEndDetectedHandler := func(event speech.RecognitionEventArgs) {
		fmt.Println("speechEndDetectedHandler", event)
		defer event.Close()
	}

	canceledHandler := func(event speech.SpeechRecognitionCanceledEventArgs) {
		fmt.Println("canceledHandler", event)
		defer event.Close()
	}

	recognizingHandler := func(event speech.SpeechRecognitionEventArgs) {
		fmt.Println("recognizingHandler", event)
		defer event.Close()
		fmt.Println("PARTIAL: ", event.Result.Text)
	}

	recognizedHandler := func(event speech.SpeechRecognitionEventArgs) {
		fmt.Println("recognizedHandler", event)
		defer event.Close()
		if event.Result.Reason == common.NoMatch {
			fmt.Println("NoMatch")
		} else {
			if event.Result.Text != "" {
				fmt.Println("FULL: ", event.Result.Text)
			}
		}
	}

	speechRecognizer.SessionStarted(sessionStartedHandler)
	speechRecognizer.SessionStopped(sessionStoppedHandler)
	speechRecognizer.SpeechStartDetected(speechStartDetectedHandler)
	speechRecognizer.SpeechEndDetected(speechEndDetectedHandler)
	speechRecognizer.Canceled(canceledHandler)
	speechRecognizer.Recognizing(recognizingHandler)
	speechRecognizer.Recognized(recognizedHandler)

	recogStartErrChan := speechRecognizer.StartContinuousRecognitionAsync()
	go func() {
		recogStartErr := <-recogStartErrChan
		fmt.Println("StartContinuousRecognitionAsync error")
		fmt.Println(recogStartErr)
		// TBD: we should probably leave AudioForkHandler here!
	}()

	go func() {
		for audioBytes := range audioStream {
			if err := stream.Write(audioBytes); err != nil {
				log.Printf("Could not send audio: %v", err)
			}
		}
	}()
}