package google

import (
	speech "cloud.google.com/go/speech/apiv1"
	"context"
	"fmt"
	"github.com/va-voice-gateway/appconfig"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"io"
	"log"
)

// Structs below are used to parse REST based Google STT configuration which are used by VAP
// These structs then needs to be transformed into respective GRPC based structs used by underlying Google APIs
type RecognitionConfig struct {
	Encoding int32                              `json:"encoding,omitempty"`
	SampleRateHertz int32                       `json:"sampleRateHertz,omitempty"`
	AudioChannelCount int32                     `json:"audioChannelCount,omitempty"`
	EnableSeparateRecognitionPerChannel bool    `json:"enableSeparateRecognitionPerChannel,omitempty"`
	LanguageCode string                         `json:"languageCode,omitempty"`
	MaxAlternatives int32                       `json:"maxAlternatives,omitempty"`
	ProfanityFilter bool                        `json:"profanityFilter,omitempty"`
	SpeechContexts []*SpeechContext             `json:"speechContexts,omitempty"`
	EnableWordTimeOffsets bool                  `json:"enableWordTimeOffsets,omitempty"`
	EnableAutomaticPunctuation bool             `json:"enableAutomaticPunctuation,omitempty"`
	DiarizationConfig *SpeakerDiarizationConfig `json:"diarizationConfig,omitempty"`
	Metadata *RecognitionMetadata               `json:"metadata,omitempty"`
	Model string                                `json:"model,omitempty"`
	UseEnhanced bool                            `json:"useEnhanced,omitempty"`
}

type RecognitionMetadata struct {
	InteractionType int32 `json:"interactionType,omitempty"`
	IndustryNaicsCodeOfAudio uint32 `json:"industryNaicsCodeOfAudio,omitempty"`
	MicrophoneDistance int32 `json:"microphoneDistance,omitempty"`
	OriginalMediaType int32 `json:"originalMediaType,omitempty"`
	RecordingDeviceType int32 `json:"recordingDeviceType,omitempty"`
	RecordingDeviceName string `json:"recordingDeviceName,omitempty"`
	OriginalMimeType string `json:"originalMimeType,omitempty"`
	AudioTopic string `json:"audioTopic,omitempty"`
}

type SpeakerDiarizationConfig struct {
	EnableSpeakerDiarization bool `json:"enableSpeakerDiarization,omitempty"`
	MinSpeakerCount int32 `json:"minSpeakerCount,omitempty"`
	MaxSpeakerCount int32 `json:"maxSpeakerCount,omitempty"`
	SpeakerTag int32 `json:"speakerTag,omitempty"`
}

// TBD: are we missing boost attribute here (same in Rust version)?
type SpeechContext struct {
	Phrases []string `json:"phrases,omitempty"`
}

// Google Speech To Text - https://cloud.google.com/speech-to-text/docs/streaming-recognize
func PerformGoogleSTT(appConfig *appconfig.AppConfig, audioStream chan []byte) {
	ctx := context.Background()
	client, err := speech.NewClient(ctx, option.WithCredentialsFile(appConfig.Temp.SttGoogleCred))
	if err != nil {
		log.Fatal(err)
	}
	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Send the initial configuration message.
	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 8000,
					LanguageCode:    "en-US",
				},
				SingleUtterance: false,
				InterimResults: true,
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	go func() {
		for audioBytes := range audioStream {
			if err := stream.Send(&speechpb.StreamingRecognizeRequest{
				StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
					AudioContent: audioBytes,
				},
			}); err != nil {
				log.Printf("Could not send audio: %v", err)
			}
		}
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Printf("StreamingRecognize EOF")
				break
			}
			if err != nil {
				log.Fatalf("Cannot stream results: %v", err)
			}
			if err := resp.Error; err != nil {
				// Workaround while the API doesn't give a more informative error.
				if err.Code == 3 || err.Code == 11 {
					log.Print("WARNING: Speech recognition request exceeded limit of 60 seconds.")
				}
				log.Fatalf("Could not recognize: %v", err)
			}
			for _, result := range resp.Results {
				fmt.Printf("Result: %+v\n", result)
			}
		}
	}()
}