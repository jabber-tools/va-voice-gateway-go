package stt

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