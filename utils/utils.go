package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/va-voice-gateway/nlpactor"
	"github.com/va-voice-gateway/logger"
	"regexp"
	"strings"
	"unicode"
)

// https://yourbasic.org/golang/regexp-cheat-sheet/
var (
	RE_DIGIT_WORDS = regexp.MustCompile(`one|two|three|four|five|six|seven|eight|nine|zero`)
	RE_SEPARATOR = regexp.MustCompile(`\s+`)
	// nice! unicode character class(\p{L}): https://stackoverflow.com/questions/30482793/golang-regexp-with-non-latin-characters
	RE_NON_ALPHANUMERIC_CHARS = regexp.MustCompile(`(?i)[^0-9a-z\p{L}\s]`)
	DIGIT_WORDS_REPLACEMENTS = []string{"1","2","3","4","5","6","7","8","9","0"}
	DIGIT_WORDS = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "zero"}
	log = logrus.New()
)

func init() {
	logger.InitLogger(log, "utils")
}

// pretty print of any structure via json marshaling with indentation
func PrettyPrint(structure interface{}) {
	b, err := json.MarshalIndent(structure, "", "  ")
	if err == nil {
		log.Debug(string(b))
	} else {
		log.Debug("BotConfigPrettyPrint error: ", err)
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
	va := nlpactor.VapActor()
	c := make(chan string)
	request := nlpactor.VapTokenRequest{Responder: c}
	va.CommandsChannel <- request
	token := <- c
	return &token
}

// go has no built in function to find position
// of element in slice
func indexOf(value string, slice []string) int {
	for p, v := range slice {
		if (v == value) {
			return p
		}
	}
	return -1
}

func isStringNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func NormalizeAWB(text string) string {
	out := text
	matches := RE_DIGIT_WORDS.FindAllStringSubmatch(out, -1)

	for _, match := range matches {
		idx := indexOf(match[0], DIGIT_WORDS)
		out = strings.Replace(out, match[0], DIGIT_WORDS_REPLACEMENTS[idx], -1 )
	}

	splits := RE_SEPARATOR.Split(out, -1)

	outSlice := []string{}

	for _, split := range splits {

		if isStringNumeric(split) {
			outSlice = append(outSlice, split)
		} else {
			outSlice = append(outSlice, fmt.Sprintf("%s ", split))
		}
	}

	return strings.Join(outSlice, "")
}

func RemoveNonAlphaNumericChars(text string) string {
	return string(RE_NON_ALPHANUMERIC_CHARS.ReplaceAll([]byte(text), []byte("")))
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}