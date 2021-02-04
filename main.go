package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/va-voice-gateway/actors"
	"github.com/va-voice-gateway/actorsvap"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/asterisk"
	"github.com/va-voice-gateway/gateway"
	"github.com/va-voice-gateway/gateway/config"
	"github.com/va-voice-gateway/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	fmt.Println("Starting Voice Gateway...")

	// load app config and caches it for global use
	appConfigPath := "c:/tmp/cfggo.toml"
	appconfig.AppConfig(&appConfigPath)
	fmt.Println("Voice Gateway config loaded")

	vapActor := actorsvap.VapActor()
	go vapActor.VapActorProcessingLoop()

	vapToken := utils.GetVapAPIToken()

	config.BotConfigs(vapToken)
	fmt.Println("Voice GW enabled Bot configs loaded")

	gatewayActor := gateway.GatewayActor()
	go gatewayActor.GatewayActorProcessingLoop()

	sttActor := actors.STTResultsActor()
	go sttActor.STTResultsActorProcessingLoop()

	ctx, cancel := context.WithCancel(context.Background())

	// termination logic
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	asterisk.Connect(ctx)
	fmt.Println("Asterisk signal stream connected!")
	go runhttp()
	<-done
	cancel()
	fmt.Println("exiting!")
}

func slashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got http request /")
}

func runhttp() {
	r := mux.NewRouter()
	r.HandleFunc("/{channelId}/{botId}/{lang}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := vars["channelId"]
		botId := vars["botId"]
		lang := vars["lang"]
		asterisk.AudioForkHandler(w, r, &channelId, &botId, &lang)
	})
	r.HandleFunc("/", slashHandler)
	fmt.Println("Listening for requests on port 8083")

	appConfig := appconfig.AppConfig(nil)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%v:%v", appConfig.Core.Host, appConfig.Core.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
