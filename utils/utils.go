package utils

import (
	"encoding/json"
	"github.com/va-voice-gateway/nlpactors"
	"log"
)

// pretty print of any structure via json marshaling with indentation
func PrettyPrint(structure interface{}) {
	b, err := json.MarshalIndent(structure, "", "  ")
	if err == nil {
		log.Println(string(b))
	} else {
		log.Println("BotConfigPrettyPrint error: ", err)
	}
}

func StructToJsonString(structure interface{}) (*string, error) {
	b, err := json.Marshal(structure)
	if err == nil {
		str := string(b)
		return &str, nil
	} else {
		return nil, err
	}
}

func GetVapAPIToken() *string {
	va := nlpactors.VapActor()
	c := make(chan string)
	request := nlpactors.VapTokenRequest{Responder: c}
	va.CommandsChannel <- request
	token := <- c
	return &token
}

