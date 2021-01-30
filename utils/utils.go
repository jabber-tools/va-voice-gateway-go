package utils

import (
	"encoding/json"
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