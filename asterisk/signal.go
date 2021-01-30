package asterisk

import (
	"context"
	"fmt"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/va-voice-gateway/appconfig"
	"log"
)

func Connect(ctx context.Context, appConfig *appconfig.AppConfig) {

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

	fmt.Println("Connected")

	info, err := cl.Asterisk().Info(nil)
	if err != nil {
		log.Fatal("Failed to get Asterisk Info", "error", err)
		return
	}

	fmt.Println("Asterisk Info", "info", info)

	fmt.Println("Starting listenAsteriskEvents")
	go listenAsteriskEvents(ctx, cl)
}

func listenAsteriskEvents(ctx context.Context, cl ari.Client) {
	subStasisStart := cl.Bus().Subscribe(nil, "StasisStart")
	subStasisEnd := cl.Bus().Subscribe(nil, "StasisEnd")
	subChannelDtmfReceived := cl.Bus().Subscribe(nil, "ChannelDtmfReceived")

	fmt.Println("listenAsteriskEvents: entering loop")
	for {
		select {
		case e := <-subStasisStart.Events():
			v := e.(*ari.StasisStart)
			fmt.Println("Got StasisStart", "channel", v.Channel.ID)
			_ = cl.Channel().Get(v.Key(ari.ChannelKey, v.Channel.ID))
		case e := <-subStasisEnd.Events():
			v := e.(*ari.StasisEnd)
			fmt.Println("Got StasisEnd", "channel", v.Channel.ID)
			_ = cl.Channel().Get(v.Key(ari.ChannelKey, v.Channel.ID))
		case e := <-subChannelDtmfReceived.Events():
			v := e.(*ari.ChannelDtmfReceived)
			fmt.Println("Got ChannelDtmfReceived", "channel", v.Channel.ID)
			fmt.Println("Digit", v.Digit)
			_ = cl.Channel().Get(v.Key(ari.ChannelKey, v.Channel.ID))
		case <-ctx.Done():
			fmt.Println("listenAsteriskEvents: leaving the loop")
			cl.Close() // disconnect from asterisk signal stream ws conn
			return
		}
	}
	fmt.Println("listenAsteriskEvents: loop left!")
}
