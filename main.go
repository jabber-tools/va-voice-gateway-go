package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/asterisk"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

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

	fmt.Println("Starting Voice Gateway...")
	appConfig, err := appconfig.LoadAppConfig("c:/tmp/cfggo.toml")
	if err != nil {
		fmt.Println("Error when loading app config")
		log.Fatal(err)
		return
	}

	fmt.Println("Voice Gateway config loaded")
	fmt.Printf("%+v\n", appConfig)

	asterisk.Connect(ctx, appConfig)
	fmt.Println("Asterisk signal stream connected!")
	go runhttp(appConfig)
	<-done
	cancel()
	fmt.Println("exiting!")
}

func slashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got http request /")
}

// placeholder for now. we will run APIs here at some moment
func runhttp(appConfig *appconfig.AppConfig) {
	r := mux.NewRouter()
	r.HandleFunc("/{channelId}/{botId}/{lang}", asterisk.AudioForkHandler)
	r.HandleFunc("/", slashHandler)
	fmt.Println("Listening for requests on port 8083")

	srv := &http.Server{
		Handler: r,
		Addr: fmt.Sprintf("%v:%v", appConfig.Core.Host, appConfig.Core.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
