package asterisk

import (
	"context"
	"fmt"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/va-voice-gateway/appconfig"
	"log"
	"net/http"
	"sync"
)

func Connect(appConfig *appconfig.AppConfig) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Connecting to Asterisk ARI")

	cl, err := native.Connect(&native.Options{
		Application:  appConfig.Asterisk.App,
		Username:     appConfig.Asterisk.Username,
		Password:     appConfig.Asterisk.Password,
		URL:          appConfig.Asterisk.AriUrl,
		WebsocketURL: appConfig.Asterisk.WSUrl,
	})
	if err != nil {
		log.Fatal("Failed to build native ARI client", "error", err)
		return
	}

	defer cl.Close()

	fmt.Println("Connected")

	info, err := cl.Asterisk().Info(nil)
	if err != nil {
		log.Fatal("Failed to get Asterisk Info", "error", err)
		return
	}

	fmt.Println("Asterisk Info", "info", info)

	fmt.Println("Starting listener app")

	go listenApp(ctx, cl, channelHandler)
	// spin up http server to prevent app from quiting. This will be replaced by AudioFork loop
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got http request /")
	}))
	fmt.Println("Listening for requests on port 9990")
	http.ListenAndServe(":9990", nil)

}

func listenApp(ctx context.Context, cl ari.Client, handler func(cl ari.Client, h *ari.ChannelHandle)) {
	sub := cl.Bus().Subscribe(nil, "StasisStart")
	end := cl.Bus().Subscribe(nil, "StasisEnd")

	for {
		select {
		case e := <-sub.Events():
			v := e.(*ari.StasisStart)
			fmt.Println("Got stasis start", "channel", v.Channel.ID)
			go handler(cl, cl.Channel().Get(v.Key(ari.ChannelKey, v.Channel.ID)))
		case <-end.Events():
			fmt.Println("Got stasis end")
		case <-ctx.Done():
			return
		}
	}
}

func channelHandler(cl ari.Client, h *ari.ChannelHandle) {
	fmt.Println("Running channel handler")

	stateChange := h.Subscribe(ari.Events.ChannelStateChange)
	defer stateChange.Cancel()

	data, err := h.Data()
	if err != nil {
		fmt.Println("Error getting data", "error", err)
		return
	}
	fmt.Println("Channel State", "state", data.State)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		fmt.Println("Waiting for channel events")

		defer wg.Done()

		for {
			select {
			case <-stateChange.Events():
				fmt.Println("Got state change request")

				data, err = h.Data()
				if err != nil {
					fmt.Println("Error getting data", "error", err)
					continue
				}
				fmt.Println("New Channel State", "state", data.State)

				if data.State == "Up" {
					stateChange.Cancel() // stop subscription to state change events
					return
				}
			}
		}
	}()

	h.Answer()

	wg.Wait()

	h.Hangup()
}