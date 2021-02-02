package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/asterisk"
	"github.com/va-voice-gateway/gateway"
	"github.com/va-voice-gateway/gateway/config"
	"github.com/va-voice-gateway/nlp"
	"github.com/va-voice-gateway/stt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	fmt.Println("Starting Voice Gateway...")
	appConfig, err := appconfig.LoadAppConfig("c:/tmp/cfggo.toml")
	if err != nil {
		fmt.Println("Error when loading app config")
		log.Fatal(err)
		return
	}

	fmt.Println("Voice Gateway config loaded")
	fmt.Printf("%+v\n", appConfig)

	vapActor := nlp.NewVapActor(appConfig.NlpVap.Username, appConfig.NlpVap.Password, appConfig.NlpVap.VapBaseUrl)
	go vapActor.VapActorProcessingLoop()

	//botCfgs, err := nlp.GetBotConfigs()
	botCfgs, err := vapActor.GetBotConfigsFromVap()
	if err != nil {
		fmt.Println("Error when loading bot configs")
		log.Fatal(err)
		return
	}

	gatewayActor := gateway.GatewayActor()
	go gatewayActor.GatewayActorProcessingLoop()

	botConfigs := config.NewBotConfigs(botCfgs)

	sttActor := stt.STTResultsActor()
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

	asterisk.Connect(ctx, appConfig)
	fmt.Println("Asterisk signal stream connected!")
	go runhttp(appConfig, &botConfigs)
	<-done
	cancel()
	fmt.Println("exiting!")
}

func slashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got http request /")
}

func runhttp(appConfig *appconfig.AppConfig, botConfigs *config.BotConfigs) {
	r := mux.NewRouter()
	r.HandleFunc("/{channelId}/{botId}/{lang}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := vars["channelId"]
		botId := vars["botId"]
		lang := vars["lang"]
		asterisk.AudioForkHandler(w, r, appConfig, &channelId, &botId, &lang, botConfigs)
	})
	r.HandleFunc("/", slashHandler)
	fmt.Println("Listening for requests on port 8083")

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%v:%v", appConfig.Core.Host, appConfig.Core.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
