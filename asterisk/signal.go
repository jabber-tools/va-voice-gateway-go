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

	gw := gateway.GatewayService()

	fmt.Println("listenAsteriskEvents: entering loop")
	for {
		select {
		case e := <-subStasisStart.Events():
			v := e.(*ari.StasisStart)
			fmt.Println("Got StasisStart", "channel", v.Channel.ID)

			go func() {
				clientId := v.Channel.ID
				botId := v.Args[0]
				lang := v.Args[1]

				channel := cl.Channel().Get(v.Key(ari.ChannelKey, v.Channel.ID))

				// TBD: load asterisk invite params
				inviteParams := make(map[string]string)
				inviteParams["foo"] = "bar"
				nlpImpl, _ := nlp.NewVAP(clientId, botId,lang, inviteParams)
				newClient := gateway.NewClient(clientId, botId, lang,inviteParams, nlpImpl)
				gw.AddClient(newClient)

				// TBD: answer call
				// TBD: nlp_tts_play (Welcome)
				channel.Answer()
			}()

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
