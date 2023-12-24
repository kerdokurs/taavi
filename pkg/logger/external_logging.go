package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kerdokurs/zlp"
	"github.com/sirupsen/logrus"
	"kerdo.dev/taavi/pkg/zulip"
)

var logStreamID string
var externalLoggingEnabled bool

func formatMessage(msg string, m M) string {
	if m == nil || len(m) == 0 {
		return msg
	}

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		log.Printf("error encoding logger.M: %v\n", err)
		return msg
	}

	fmtBuf := bytes.Buffer{}
	if err := json.Indent(&fmtBuf, buf.Bytes(), "", "  "); err != nil {
		log.Printf("error indenting logger.M: %v\n", err)
		return msg
	}

	return fmt.Sprintf("%s\n\n```\n%s\n```", msg, fmtBuf.String())
}

func logExternal(level logrus.Level, msg string, m M) {
	var logTopicID string
	switch level {
	case logrus.ErrorLevel:
		logTopicID = "error"
	case logrus.FatalLevel:
		logTopicID = "fatal"
	case logrus.InfoLevel:
		fallthrough
	default:
		logTopicID = "info"
	}

	if externalLoggingEnabled {
		msg := zlp.Message{
			Stream:  logStreamID,
			Topic:   logTopicID,
			Content: formatMessage(msg, m),
		}
		zulip.Client.Message(&msg)
	}
}
