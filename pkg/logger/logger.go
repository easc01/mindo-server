package logger

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetLevel(logrus.TraceLevel)
	Log.SetReportCaller(true) // Enable built-in caller tracking
	Log.SetFormatter(&CustomFormatter{})
}

// CustomFormatter formats log entries with caller info
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var file string
	var line int

	if entry.Caller != nil {
		// If using Go modules, this cleans up the file path
		file = trimGoPath(entry.Caller.File)
		line = entry.Caller.Line
	}

	logLine := fmt.Sprintf("[%s] [%s] [%s]:[%d] - %s\n",
		entry.Time.Format("2006-01-02 15:04:05"),
		strings.ToUpper(entry.Level.String()),
		file,
		line,
		entry.Message,
	)

	return []byte(logLine), nil
}

func trimGoPath(fullPath string) string {
	// Trim GOPATH or module path if needed
	parts := strings.Split(fullPath, "/mindo-server/")
	if len(parts) > 1 {
		return "mindo-server/" + parts[1]
	}
	return fullPath
}

// Wrappers to avoid leaking caller info from here
func Info(args ...interface{})  { Log.Info(args...) }
func Warn(args ...interface{})  { Log.Warn(args...) }
func Error(args ...interface{}) { Log.Error(args...) }
func Debug(args ...interface{}) { Log.Debug(args...) }
