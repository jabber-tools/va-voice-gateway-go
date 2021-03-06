package google

import (
	speech "cloud.google.com/go/speech/apiv1"
	"context"
	"fmt"
	"github.com/va-voice-gateway/gateway"
	"github.com/va-voice-gateway/logger"
	"github.com/va-voice-gateway/sttactor"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/gateway/config"
	"github.com/va-voice-gateway/utils"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"io"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	logger.InitLogger(log, "main")
}

// truly quick & dirty, see rust based implementation for proper stuff
func IntoGrpc(rc *config.RecognitionConfig, lang *string) *speechpb.RecognitionConfig {
	var encoding speechpb.RecognitionConfig_AudioEncoding
	switch rc.Encoding {
	case 0:
		encoding = speechpb.RecognitionConfig_LINEAR16 // default
	case 1:
		encoding = speechpb.RecognitionConfig_LINEAR16
	case 2:
		encoding = speechpb.RecognitionConfig_FLAC
	case 3:
		encoding = speechpb.RecognitionConfig_MULAW
	case 4:
		encoding = speechpb.RecognitionConfig_AMR
	case 5:
		encoding = speechpb.RecognitionConfig_AMR_WB
	case 6:
		encoding = speechpb.RecognitionConfig_OGG_OPUS
	case 7:
		encoding = speechpb.RecognitionConfig_SPEEX_WITH_HEADER_BYTE
	}

	var languageCode *string
	if rc.LanguageCode != "" {
		languageCode = &rc.LanguageCode
	} else {
		languageCode = lang
	}

	var ctxs []*speechpb.SpeechContext = make([]*speechpb.SpeechContext, len(rc.SpeechContexts))

	for idx, ctx := range rc.SpeechContexts {
		ctxs[idx] = &speechpb.SpeechContext{
			Phrases: ctx.Phrases,
		}
	}

	var it speechpb.RecognitionMetadata_InteractionType = speechpb.RecognitionMetadata_INTERACTION_TYPE_UNSPECIFIED
	if rc.Metadata != nil {
		switch rc.Metadata.InteractionType {
		case 0:
			it = speechpb.RecognitionMetadata_INTERACTION_TYPE_UNSPECIFIED
		case 1:
			it = speechpb.RecognitionMetadata_DISCUSSION
		case 2:
			it = speechpb.RecognitionMetadata_PRESENTATION
		case 3:
			it = speechpb.RecognitionMetadata_PHONE_CALL
		case 4:
			it = speechpb.RecognitionMetadata_VOICEMAIL
		case 5:
			it = speechpb.RecognitionMetadata_PROFESSIONALLY_PRODUCED
		case 6:
			it = speechpb.RecognitionMetadata_VOICE_SEARCH
		case 7:
			it = speechpb.RecognitionMetadata_VOICE_COMMAND
		case 8:
			it = speechpb.RecognitionMetadata_DICTATION
		}
	}

	var md speechpb.RecognitionMetadata_MicrophoneDistance = speechpb.RecognitionMetadata_MICROPHONE_DISTANCE_UNSPECIFIED
	if rc.Metadata != nil {
		switch rc.Metadata.MicrophoneDistance {
		case 0:
			md = speechpb.RecognitionMetadata_MICROPHONE_DISTANCE_UNSPECIFIED
		case 1:
			md = speechpb.RecognitionMetadata_NEARFIELD
		case 2:
			md = speechpb.RecognitionMetadata_MIDFIELD
		case 3:
			md = speechpb.RecognitionMetadata_FARFIELD
		}
	}

	var omt speechpb.RecognitionMetadata_OriginalMediaType = speechpb.RecognitionMetadata_ORIGINAL_MEDIA_TYPE_UNSPECIFIED
	if rc.Metadata != nil {
		switch rc.Metadata.OriginalMediaType {
		case 0:
			omt = speechpb.RecognitionMetadata_ORIGINAL_MEDIA_TYPE_UNSPECIFIED
		case 1:
			omt = speechpb.RecognitionMetadata_AUDIO
		case 2:
			omt = speechpb.RecognitionMetadata_VIDEO
		}
	}

	var rdt speechpb.RecognitionMetadata_RecordingDeviceType = speechpb.RecognitionMetadata_RECORDING_DEVICE_TYPE_UNSPECIFIED
	if rc.Metadata != nil {
		switch rc.Metadata.RecordingDeviceType {
		case 0:
			rdt = speechpb.RecognitionMetadata_RECORDING_DEVICE_TYPE_UNSPECIFIED
		case 1:
			rdt = speechpb.RecognitionMetadata_SMARTPHONE
		case 2:
			rdt = speechpb.RecognitionMetadata_PC
		case 3:
			rdt = speechpb.RecognitionMetadata_PHONE_LINE
		case 4:
			rdt = speechpb.RecognitionMetadata_VEHICLE
		case 5:
			rdt = speechpb.RecognitionMetadata_OTHER_OUTDOOR_DEVICE
		case 6:
			rdt = speechpb.RecognitionMetadata_OTHER_INDOOR_DEVICE
		}
	}

	var dc *speechpb.SpeakerDiarizationConfig
	if rc.DiarizationConfig != nil {
		dc = &speechpb.SpeakerDiarizationConfig{
			EnableSpeakerDiarization: rc.DiarizationConfig.EnableSpeakerDiarization,
			MinSpeakerCount:          rc.DiarizationConfig.MinSpeakerCount,
			MaxSpeakerCount:          rc.DiarizationConfig.MaxSpeakerCount,
		}
	} else {
		dc = nil
	}

	var rmt *speechpb.RecognitionMetadata
	if rc.Metadata != nil {
		rmt = &speechpb.RecognitionMetadata{
			InteractionType:          it,
			IndustryNaicsCodeOfAudio: rc.Metadata.IndustryNaicsCodeOfAudio,
			MicrophoneDistance:       md,
			OriginalMediaType:        omt,
			RecordingDeviceType:      rdt,
			RecordingDeviceName:      rc.Metadata.RecordingDeviceName,
			OriginalMimeType:         rc.Metadata.OriginalMimeType,
			AudioTopic:               rc.Metadata.AudioTopic,
		}
	} else {
		rmt = nil
	}

	var sampleRateHertz int32
	if rc.SampleRateHertz != 0 {
		sampleRateHertz = rc.SampleRateHertz
	} else {
		sampleRateHertz = 8000
	}

	var audioChannelCount int32
	if rc.AudioChannelCount != 0 {
		audioChannelCount = rc.AudioChannelCount
	} else {
		audioChannelCount = 1
	}

	rcout := speechpb.RecognitionConfig{
		Encoding:                            encoding,
		SampleRateHertz:                     sampleRateHertz,
		AudioChannelCount:                   audioChannelCount,
		EnableSeparateRecognitionPerChannel: rc.EnableSeparateRecognitionPerChannel,
		LanguageCode:                        *languageCode,
		MaxAlternatives:                     rc.MaxAlternatives,
		ProfanityFilter:                     rc.ProfanityFilter,
		SpeechContexts:                      ctxs,
		EnableWordTimeOffsets:               rc.EnableWordTimeOffsets,
		EnableAutomaticPunctuation:          rc.EnableAutomaticPunctuation,
		DiarizationConfig:                   dc,
		Metadata:                            rmt,
		Model:                               rc.Model,
		UseEnhanced:                         rc.UseEnhanced,
	}

	// utils.PrettyPrint(rcout)

	return &rcout
}

// Google Speech To Text - https://cloud.google.com/speech-to-text/docs/streaming-recognize
func PerformGoogleSTT(audioStream *chan []byte, recCfg *config.RecognitionConfig, botId *string, channelId *string, lang *string, signalToAudioFork *chan int) {
	log.Infof("PerformGoogleSTT called for channel %v\n", *channelId)
	ctx := context.Background()
	_ = appconfig.AppConfig(nil) // not needed for now

	botConfigs := config.BotConfigs(nil)

	credStr, err := utils.StructToJsonString(botConfigs.GetSTTGoogleCred(botId))
	if err != nil {
		log.Error("Unable to retrieve Google STT Credentials")
	}

	credBytes := []byte(*credStr)
	client, err := speech.NewClient(ctx, option.WithCredentialsJSON(credBytes))
	if err != nil {
		log.Error(err)
	}
	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Error(err)
	}
	// Send the initial configuration message.
	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config:          IntoGrpc(recCfg, lang),
				SingleUtterance: false,
				InterimResults:  true,
			},
		},
	}); err != nil {
		log.Error(err)
	}

	go func() {
		gw := gateway.GatewayService()
		for audioBytes := range *audioStream {
			if gw.GetDoSTT(channelId) == true {
				if err := stream.Send(&speechpb.StreamingRecognizeRequest{
					StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
						AudioContent: audioBytes,
					},
				}); err != nil {
					log.Errorf("Could not send audio: %v", err)
				}
			}
		}
		log.Infof("PerformGoogleSTT go loop #1 left: %v", *channelId)
	}()

	go func(channelId *string) {
		for {
			resp, err := stream.Recv()
			// log.Printf(">>>resp %v\n", resp)
			// log.Printf(">>>err %v\n", err)
			/*
			2021/02/09 21:28:23 Got StasisEnd channel 1612902439.8
			2021/02/09 21:28:33 >>>resp error:{code:11  message:"Audio Timeout Error: Long duration elapsed without audio. Audio should be sent close to real time."}
			2021/02/09 21:28:33 >>>err <nil>
			2021/02/09 21:28:33 WARNING: Speech recognition request exceeded limit of 60 seconds.
			2021/02/09 21:28:33 STTResultsActorProcessingLoop.CommandErrorResult  {1612902439.8 Could not recognize: code:11  message:"Audio Timeout Error: Long duration elapsed without audio. Audio should be sent close to real time."
			}
			2021/02/09 21:28:33 >>>resp <nil>
			2021/02/09 21:28:33 >>>err rpc error: code = OutOfRange desc = Audio Timeout Error: Long duration elapsed without audio. Audio should be sent close to real time.
			panic: runtime error: invalid memory address or nil pointer dereference
			[signal 0xc0000005 code=0x0 addr=0x28 pc=0x916cee]

			goroutine 41 [running]:
			github.com/va-voice-gateway/stt/google.PerformGoogleSTT.func2(0xbc1b40, 0xc00049c630, 0xc0002fa740)
				C:/Users/abezecny/GoLandProjects/va-voice-gateway-go/stt/google/google.go:273 +0x36e
			created by github.com/va-voice-gateway/stt/google.PerformGoogleSTT
				C:/Users/abezecny/GoLandProjects/va-voice-gateway-go/stt/google/google.go:235 +0x471

			-> once we retrieve resp with err code 11 we need to terminate the loop, next call will never succeed
			*/

			if err == io.EOF {
				log.Info("StreamingRecognize EOF")
				break
			}
			if err != nil {
				// log.Printf("Cannot stream results: %v\n", err)
				sttactor.STTResultsActor().CommandsChannel <- sttactor.CommandErrorResult{
					ChannelId: *channelId,
					Error: err,
				}
			}
			if err := resp.Error; err != nil {
				// Workaround while the API doesn't give a more informative error.
				if err.Code == 3 || err.Code == 11 {
					log.Warn("WARNING: Speech recognition request exceeded limit of 60 seconds.")
					sttactor.STTResultsActor().CommandsChannel <- sttactor.CommandErrorResult{
						ChannelId: *channelId,
						Error: fmt.Errorf("%v\n", err),
					}
					log.Debugf("sending signalToAudioFork = 1: %v", *channelId)
					*signalToAudioFork <- 1 // value 1 indicates audiofork should spin up another PerformGoogleSTT go routine to recover
					log.Debugf("sent signalToAudioFork = 1: %v", *channelId)
					break // no point to do next iteration we need to call PerformGoogleSTT again
				}
				sttactor.STTResultsActor().CommandsChannel <- sttactor.CommandErrorResult{
					ChannelId: *channelId,
					Error: fmt.Errorf("Could not recognize: %v\n", err),
				}
			}
			for _, result := range resp.Results {
				if result.IsFinal == true {
					sttactor.STTResultsActor().CommandsChannel <- sttactor.CommandFinalResult{
						ChannelId: *channelId,
						Text: result.Alternatives[0].Transcript,
					}
				} else {
					sttactor.STTResultsActor().CommandsChannel <- sttactor.CommandPartialResult{
						ChannelId: *channelId,
						Text: result.Alternatives[0].Transcript,
					}
				}
			}
		}
		log.Infof("PerformGoogleSTT go loop #2 left: %v", *channelId)
	}(channelId)
}
