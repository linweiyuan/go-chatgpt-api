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
	green  = Ansi("\033[1;32m%s\033[0m")
	yellow = Ansi("\033[1;33m%s\033[0m")
	red    = Ansi("\033[1;31m%s\033[0m")
)

func Info(msg string) {
	logrus.Info(green(msg))
}

func Warn(msg string) {
	logrus.Warn(yellow(msg))
}

func Error(msg string) {
	logrus.Error(red(msg))
}
