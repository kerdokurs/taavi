package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type M = logrus.Fields

func Init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	logStreamID = "TEAVITUSED (TEST)"
	externalLoggingEnabled = os.Getenv("LOGGING") == "true"
}

func Infow(msg string, m M) {
	logrus.WithFields(m).Info(msg)
	logExternal(logrus.InfoLevel, msg, m)
}

func Errorw(msg string, m M) {
	logrus.WithFields(m).Error(msg)
	logExternal(logrus.ErrorLevel, msg, m)
}

func Fatalw(msg string, m M) {
	logrus.WithFields(m).Fatal(msg)
	logExternal(logrus.FatalLevel, msg, m)
	os.Exit(1)
}
