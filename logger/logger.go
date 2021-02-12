package logger

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

// must be duplicated from utils to prevent import cycle :(
func structToJsonString(structure interface{}) (*string, error) {
	b, err := json.Marshal(structure)
	if err == nil {
		str := string(b)
		return &str, nil
	} else {
		return nil, err
	}
}

type customLogFormatter struct {}

func (f *customLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	Level, _ := entry.Level.MarshalText()
	LevelStr := fmt.Sprintf("[%s]", Level)

	// TimeStr := entry.Time.Format("2006-01-02 15:04:05.999999999 -0700 MST")
	TimeStr := fmt.Sprintf("[%s]", entry.Time.Format("2006-01-02 15:04:05.999 -0700 MST")) // ms precision is OK

	ContextData, _ := structToJsonString(entry.Data)

	ContextDataStr := ""
	if *ContextData != "{}" {
		ContextDataStr = fmt.Sprintf(" %s",*ContextData)
	}

	// github.com/va-voice-gateway/appconfig -> appconfig
	caller := fmt.Sprintf("[%s]", strings.ReplaceAll(entry.Caller.Function, "github.com/va-voice-gateway/", ""))
	msg := fmt.Sprintf("%s %s %s %s%s\n", TimeStr, string(LevelStr), caller, entry.Message, ContextDataStr)

	return []byte(msg), nil
}

func InitLogger(logger *logrus.Logger) {
	logger.SetReportCaller(true)
	logger.SetFormatter(new(customLogFormatter))
}
