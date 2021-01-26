module github.com/va-voice-gateway

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/CyCoreSystems/ari/v5 v5.1.2
	github.com/Microsoft/cognitive-services-speech-sdk-go v1.14.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
)

// replacement with forked version which does not require client
// websocket frames to be masked as per RFC6455 requirements
replace github.com/gorilla/websocket v1.4.2 => github.com/adambezecny/websocket v1.4.3
