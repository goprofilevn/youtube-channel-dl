package logger

import (
	"fmt"

	"github.com/golang-module/carbon/v2"
	"github.com/gookit/color"
)

type LogLevel string

const (
	LevelError   LogLevel = "error"
	LevelWarn    LogLevel = "warn"
	LevelDebug   LogLevel = "debug"
	LevelSuccess LogLevel = "success"
	LevelInfo    LogLevel = "info"
)

func logger(level LogLevel, message any) {
	status := ""
	switch level {
	case LevelError:
		status = color.Red.Render("[error]")
		break
	case LevelDebug:
		status = color.Blue.Render("[debug]")
		break
	case LevelSuccess:
		status = color.Green.Render("[success]")
		break
	case LevelWarn:
		status = color.Yellow.Render("[warn]")
		break
	default:
		status = color.Cyan.Render("[info]")
		break
	}
	timeNow := carbon.Now().Format("d/m/Y H:i:s")
	fmt.Printf("%s %s %s\n", color.Gray.Render("["+timeNow+"]"), status, message)
}

func Error(message any) {
	logger(LevelError, message)
}
func Errorf(format string, args ...any) {
	logger(LevelError, fmt.Sprintf(format, args...))
}
func Warn(message any) {
	logger(LevelWarn, message)
}
func Warnf(format string, args ...any) {
	logger(LevelWarn, fmt.Sprintf(format, args...))
}
func Debug(message any) {
	logger(LevelDebug, message)
}
func Debugf(format string, args ...any) {
	logger(LevelDebug, fmt.Sprintf(format, args...))
}
func Success(message any) {
	logger(LevelSuccess, message)
}
func Successf(format string, args ...any) {
	logger(LevelSuccess, fmt.Sprintf(format, args...))
}
func Info(message any) {
	logger(LevelInfo, message)
}
func Infof(format string, args ...any) {
	logger(LevelInfo, fmt.Sprintf(format, args...))
}

func Print(level LogLevel, status string, message any) {
	switch level {
	case LevelError:
		status = color.Red.Render("[" + status + "]")
		break
	case LevelDebug:
		status = color.Blue.Render("[" + status + "]")
		break
	case LevelSuccess:
		status = color.Green.Render("[" + status + "]")
		break
	case LevelWarn:
		status = color.Yellow.Render("[" + status + "]")
		break
	default:
		status = color.Cyan.Render("[" + status + "]")
		break
	}
	timeNow := carbon.Now().Format("d/m/Y H:i:s")
	fmt.Printf("%s %s %s\n", color.Gray.Render("["+timeNow+"]"), status, message)
}
