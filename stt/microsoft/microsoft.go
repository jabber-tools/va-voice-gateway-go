// +build linux

package microsoft

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"google.golang.org/appengine/log"
)

func PerformMicrosoftSTT(audioStream chan []byte, botId *string, channelId *string, lang *string) {
	audioFormat, err := audio.GetWaveFormatPCM(8000, 16, 1)
	if err != nil {
		log.Error("GetWaveFormatPCM error: ", err)
		return
	}
	defer audioFormat.Close()

	stream, err := audio.CreatePushAudioInputStreamFromFormat(audioFormat)
	if err != nil {
		log.Error("CreatePushAudioInputStreamFromFormat error: ", err)
		return
	}
	defer stream.Close()

	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
	if err != nil {
		log.Error("NewAudioConfigFromStreamInput error: ", err)
		return
	}
	defer audioConfig.Close()

	speechConfig, err := speech.NewSpeechConfigFromSubscription(appConfig.Temp.SttMsSubKey, appConfig.Temp.SttMsRegion)
	if err != nil {
		log.Error("NewAudioConfigFromStreamInput error: ", err)
		return
	}
	defer speechConfig.Close()

	speechRecognizer, err := speech.NewSpeechRecognizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		log.Error("NewSpeechRecognizerFromConfig error: ", err)
		return
	}
	defer speechRecognizer.Close()

	sessionStartedHandler := func(event speech.SessionEventArgs) {
		log.Debug("sessionStartedHandler", event)
		defer event.Close()
	}

	sessionStoppedHandler := func(event speech.SessionEventArgs) {
		log.Debug("sessionStoppedHandler", event)
		defer event.Close()
	}

	speechStartDetectedHandler := func(event speech.RecognitionEventArgs) {
		log.Debug("speechStartDetectedHandler", event)
		defer event.Close()
	}

	speechEndDetectedHandler := func(event speech.RecognitionEventArgs) {
		log.Debug("speechEndDetectedHandler", event)
		defer event.Close()
	}

	canceledHandler := func(event speech.SpeechRecognitionCanceledEventArgs) {
		log.Debug("canceledHandler", event)
		defer event.Close()
	}

	recognizingHandler := func(event speech.SpeechRecognitionEventArgs) {
		log.Debug("recognizingHandler", event)
		defer event.Close()
		log.Debug("PARTIAL: ", event.Result.Text)
	}

	recognizedHandler := func(event speech.SpeechRecognitionEventArgs) {
		log.Debug("recognizedHandler", event)
		defer event.Close()
		if event.Result.Reason == common.NoMatch {
			log.Debug("NoMatch")
		} else {
			if event.Result.Text != "" {
				log.Debug("FULL: ", event.Result.Text)
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
		log.Error("StartContinuousRecognitionAsync error")
		log.Error(recogStartErr)
		// TBD: we should probably leave AudioForkHandler here!
	}()

	go func() {
		for audioBytes := range audioStream {
			if err := stream.Write(audioBytes); err != nil {
				log.Errorf("Could not send audio: %v", err)
			}
		}
	}()
}