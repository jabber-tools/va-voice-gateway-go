package logger

// see https://pkg.go.dev/github.com/sirupsen/logrus#readme-logrus
// https://github.com/Sirupsen/logrus

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
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
	Level, err := entry.Level.MarshalText()
	if err != nil {
		return nil, err
	}
	LevelStr := fmt.Sprintf("[%s]", Level)

	// TimeStr := entry.Time.Format("2006-01-02 15:04:05.999999999 -0700 MST")
	TimeStr := fmt.Sprintf("[%s]", entry.Time.Format("2006-01-02 15:04:05.999 -0700 MST")) // ms precision is OK

	ContextData, err := structToJsonString(entry.Data)

	if err != nil {
		return nil, err
	}

	ContextDataStr := ""
	if *ContextData != "{}" {
		ContextDataStr = fmt.Sprintf(" %s",*ContextData)
	}

	// github.com/va-voice-gateway/appconfig -> appconfig
	caller := fmt.Sprintf("[%s]", strings.ReplaceAll(entry.Caller.Function, "github.com/va-voice-gateway/", ""))
	msg := fmt.Sprintf("%s %s %s %s%s\n", TimeStr, string(LevelStr), caller, entry.Message, ContextDataStr)

	return []byte(msg), nil
}

func InitLogger(logger *logrus.Logger, PackageName string) {

	packageLoggingLevel := os.Getenv(fmt.Sprintf("loglevel_%s", PackageName))

	var  logLevel logrus.Level
	switch packageLoggingLevel {
		case "panic":
			logLevel = logrus.PanicLevel
			break
		case "fatal":
			logLevel = logrus.FatalLevel
			break
		case "error":
			logLevel = logrus.ErrorLevel
			break
		case "warn":
			logLevel = logrus.WarnLevel
			break
		case "info":
			logLevel = logrus.InfoLevel
			break
		case "debug":
			logLevel = logrus.DebugLevel
			break
		case "trace":
			logLevel = logrus.TraceLevel
			break
		default:
			logLevel = logrus.DebugLevel
	}

	logger.SetReportCaller(true)
	logger.SetFormatter(new(customLogFormatter))
	logger.SetLevel(logLevel)
}
