package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
}

func Ansi(colorString string) func(...interface{}) string {
	return func(args ...interface{}) string {
		return fmt.Sprintf(colorString, fmt.Sprint(args...))
	}
}

var (
	Green  = Ansi("\033[1;32m%s\033[0m")
	Yellow = Ansi("\033[1;33m%s\033[0m")
	Red    = Ansi("\033[1;31m%s\033[0m")
)

func Info(msg string) {
	logrus.Info(Green(msg))
}

func Warn(msg string) {
	logrus.Warn(Yellow(msg))
}

func Error(msg string) {
	logrus.Error(Red(msg))
}
