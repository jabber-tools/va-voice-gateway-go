package asterisk

import (
	"context"
	"fmt"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/va-voice-gateway/appconfig"
	"github.com/va-voice-gateway/asteriskclient"
	"github.com/va-voice-gateway/gateway"
	"github.com/va-voice-gateway/nlp"
	"github.com/va-voice-gateway/utils"
	"log"
	"strings"
)

var (
	DIGITS = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
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
	subChannelHangupRequest := cl.Bus().Subscribe(nil, "ChannelHangupRequest")
	subChannelTalkingStarted := cl.Bus().Subscribe(nil, "ChannelTalkingStarted")
	subChannelTalkingFinished := cl.Bus().Subscribe(nil, "ChannelTalkingFinished")
	subChannelDestroyed := cl.Bus().Subscribe(nil, "ChannelDestroyed")
	subPlaybackFinished := cl.Bus().Subscribe(nil, "PlaybackFinished")
	subPlaybackStarted	 := cl.Bus().Subscribe(nil, "PlaybackStarted")

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
			case e := <-subChannelHangupRequest.Events():
				event := e.(*ari.ChannelHangupRequest)
				go handlerChannelHangupRequest(event, cl)
			case e := <-subChannelTalkingStarted.Events():
				event := e.(*ari.ChannelTalkingStarted)
				go handlerChannelTalkingStarted(event, cl)
			case e := <-subChannelTalkingFinished.Events():
				event := e.(*ari.ChannelTalkingFinished)
				go handlerChannelTalkingFinished(event, cl)
			case e := <-subChannelDestroyed.Events():
				event := e.(*ari.ChannelDestroyed)
				go handlerChannelDestroyed(event, cl)
			case e := <-subPlaybackFinished.Events():
				event := e.(*ari.PlaybackFinished)
				go handlerPlaybackFinished(event, cl)
			case e := <-subPlaybackStarted.Events():
				event := e.(*ari.PlaybackStarted)
				go handlerPlaybackStarted(event, cl)
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

	inviteParams := make(map[string]string)
	// does not work
	//varVARS, err := cl.Asterisk().Variables().Get(ari.NewKey(ari.VariableKey, "VARS"))
	varVARS, err := channel.GetVariable("VARS")
	if err != nil {
		log.Printf("error when loading asterisk variables")
	} else {
		log.Printf("asterisk variables " + varVARS)
		splits := strings.Split(varVARS, ",")
		for _, varXname := range splits {
			log.Printf("Loading asterisk variable " + varXname)
			varXval, err := channel.GetVariable(varXname)
			if err != nil {
				log.Printf("Error when loading asterisk variable " + varXname)
				continue
			} else {
				inviteParams[varXname] = varXval
			}
		}
	}

	nlpImpl, _ := nlp.NewVAP(clientId, botId,lang, inviteParams)
	newClient := gateway.NewClient(clientId, botId, lang,inviteParams, nlpImpl)
	gw.AddClient(newClient)

	channel.Answer()

	go asteriskclient.Nlp_tts_play(&clientId, &botId, &lang, nlp.NLPRequest{
		Text: nil,
		Event: &nlp.NLPRequestEvent{
			Name: "Welcome",
		},
	})
}

func handlerStasisEnd(event *ari.StasisEnd, cl ari.Client) {
	log.Println("Got StasisEnd", "channel", event.Channel.ID)
}

func handlerChannelDtmfReceived(event *ari.ChannelDtmfReceived, cl ari.Client) {
	log.Println("Got ChannelDtmfReceived", "channel", event.Channel.ID, "Digit", event.Digit)

	gw := gateway.GatewayService()
	clientId := event.Channel.ID

	if event.Digit == "#" {
		dtmf := gw.GetDtmf(&clientId)
		log.Println("final dtmf ", dtmf)
		gw.ResetDtmf(&clientId)

		botId, lang := gw.GetBotIdLang(&clientId)
		if botId!=nil && lang != nil {
			go asteriskclient.Nlp_tts_play(&clientId, botId, lang, nlp.NLPRequest{
				Text: &nlp.NLPRequestText{
					Text: dtmf,
				},
				Event: nil,
			})
		}
	} else {
		if utils.Contains(DIGITS, event.Digit) {
			gw.AddDtmf(&clientId, event.Digit)
			log.Println("Adding dtmf ",  event.Digit, event.Channel.ID)
		}
	}
}

func handlerChannelHangupRequest(event *ari.ChannelHangupRequest, cl ari.Client) {
	log.Println("Got ChannelHangupRequest", "channel", event.Channel.ID)
}

func handlerChannelTalkingStarted(event *ari.ChannelTalkingStarted, cl ari.Client) {
	log.Println("Got ChannelTalkingStarted", "channel", event.Channel.ID)
}

func handlerChannelTalkingFinished(event *ari.ChannelTalkingFinished, cl ari.Client) {
	log.Println("Got ChannelTalkingFinished", "channel", event.Channel.ID)
}

func handlerChannelDestroyed(event *ari.ChannelDestroyed, cl ari.Client) {
	log.Println("Got ChannelDestroyed", "channel", event.Channel.ID)
}

func handlerPlaybackFinished(event *ari.PlaybackFinished, cl ari.Client) {
	log.Println("Got PlaybackFinished")
}

func handlerPlaybackStarted(event *ari.PlaybackStarted, cl ari.Client) {
	log.Println("Got PlaybackStarted")
}
