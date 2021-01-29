module github.com/va-voice-gateway

go 1.15

require (
	cloud.google.com/go v0.75.0
	github.com/BurntSushi/toml v0.3.1
	github.com/CyCoreSystems/ari/v5 v5.1.2
	github.com/Microsoft/cognitive-services-speech-sdk-go v1.14.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/oauth2 v0.0.0-20210126194326-f9ce19ea3013 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/api v0.37.0
	google.golang.org/genproto v0.0.0-20210126160654-44e461bb6506
	google.golang.org/grpc v1.35.0 // indirect
)

// replacement with forked version which does not require client
// websocket frames to be masked as per RFC6455 requirements
replace github.com/gorilla/websocket v1.4.2 => github.com/adambezecny/websocket v1.4.3
