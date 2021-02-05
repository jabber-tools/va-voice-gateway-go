package asterisk

import (
	"context"
	"fmt"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/gateway"
	"github.com/va-voice-gateway/nlp"
	"log"
)

func Connect(ctx context.Context) *ari.Client {

	fmt.Println("Connecting to Asterisk ARI")

	appConfig := appconfig.AppConfig(nil)

	cl, err := native.Connect(&native.Options{
		Application:  appConfig.Asterisk.App,
		Username:     appConfig.Asterisk.Username,
		Password:     appConfig.Asterisk.Password,
		URL:          appConfig.Asterisk.AriUrl,
		WebsocketURL: appConfig.Asterisk.WSUrl,
	})
	if err != nil {
		log.Fatal("Failed to build native ARI client", "error", err)
		return nil
	}

	fmt.Println("Connected")

	info, err := cl.Asterisk().Info(nil)
	if err != nil {
		log.Fatal("Failed to get Asterisk Info", "error", err)
		return nil
	}

	fmt.Println("Asterisk Info", "info", info)

	fmt.Println("Starting listenAsteriskEvents")
	go listenAsteriskEvents(ctx, cl)

	return &cl
}

func listenAsteriskEvents(ctx context.Context, cl ari.Client) {
	subStasisStart := cl.Bus().Subscribe(nil, "StasisStart")
	subStasisEnd := cl.Bus().Subscribe(nil, "StasisEnd")
	subChannelDtmfReceived := cl.Bus().Subscribe(nil, "ChannelDtmfReceived")

	fmt.Println("listenAsteriskEvents: entering loop")
	for {
		select {
			case e := <-subStasisStart.Events():
				event := e.(*ari.StasisStart)
				go handlerStasisStart(event, cl)
			case e := <-subStasisEnd.Events():
				event := e.(*ari.StasisEnd)
				go handlerStasisEnd(event, cl)
			case e := <-subChannelDtmfReceived.Events():
				event := e.(*ari.ChannelDtmfReceived)
				go handlerChannelDtmfReceived(event, cl)
			case <-ctx.Done():
				fmt.Println("listenAsteriskEvents: leaving the loop")
				cl.Close() // disconnect from asterisk signal stream ws conn
				return
		}
	}
	fmt.Println("listenAsteriskEvents: loop left!")
}

func handlerStasisStart(event *ari.StasisStart, cl ari.Client) {
	fmt.Println("Got StasisStart", "channel", event.Channel.ID)

	gw := gateway.GatewayService()

	clientId := event.Channel.ID
	botId := event.Args[0]
	lang := event.Args[1]

	channel := cl.Channel().Get(event.Key(ari.ChannelKey, event.Channel.ID))

	// TBD: load asterisk invite params
	inviteParams := make(map[string]string)
	inviteParams["foo"] = "bar"
	nlpImpl, _ := nlp.NewVAP(clientId, botId,lang, inviteParams)
	newClient := gateway.NewClient(clientId, botId, lang,inviteParams, nlpImpl)
	gw.AddClient(newClient)

	// TBD: answer call
	// TBD: nlp_tts_play (Welcome)
	channel.Answer()
}

func handlerStasisEnd(event *ari.StasisEnd, cl ari.Client) {
	fmt.Println("Got StasisEnd", "channel", event.Channel.ID)
}

func handlerChannelDtmfReceived(event *ari.ChannelDtmfReceived, cl ari.Client) {
	fmt.Println("Got ChannelDtmfReceived", "channel", event.Channel.ID)
	fmt.Println("Digit", event.Digit)
}
