package main

import (
	"context"
	"fmt"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/asterisk"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	go runhttp()
	<-done
	cancel()
	fmt.Println("exiting!")
}

// placeholder for now. we will run APIs here at some moment
func runhttp() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got http request /")
	}))
	fmt.Println("Listening for requests on port 9990")
	http.ListenAndServe(":9990", nil)
}
