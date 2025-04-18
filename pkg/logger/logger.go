package logger

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	Log.SetLevel(logrus.TraceLevel)

	Log.SetFormatter(&CustomFormatter{})
}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	_, file, line, _ := runtime.Caller(6)
	logLine := fmt.Sprintf("[%s] [%s] [%s]:[%d] - %s\n", entry.Time.Format("2006-01-02 15:04:05"), entry.Level, file, line, entry.Message)
	return []byte(logLine), nil
}
