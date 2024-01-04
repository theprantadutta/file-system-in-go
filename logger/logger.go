package logger

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"testing_bubble_tea/colors"
	"time"
)

type logLevel int

const (
	DebugLevel logLevel = iota
	InfoLevel
	ErrorLevel
	SuccessLevel
	WarningLevel
)

func log(level logLevel, msg string) {

	timestamp := time.Now().Format("02-Jan-2006 03:04:05PM")

	logColor, logStr := getLogColor(level)

	var levelGloss = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors.WhiteTextColor)).
		Background(lipgloss.Color(logColor)).
		Width(10).
		AlignHorizontal(lipgloss.Center).
		PaddingLeft(1).
		PaddingRight(1)

	fmt.Println(timestamp, levelGloss.Render(logStr), msg)
}

func Success(args ...interface{}) {
	msg := fmt.Sprint(args...)
	log(SuccessLevel, msg)
}

func Info(args ...interface{}) {
	msg := fmt.Sprint(args...)
	log(InfoLevel, msg)
}

func Error(args ...interface{}) {
	msg := fmt.Sprint(args...)
	log(ErrorLevel, msg)
}

func Warning(args ...interface{}) {
	msg := fmt.Sprint(args...)
	log(WarningLevel, msg)
}

func Debug(args ...interface{}) {
	msg := fmt.Sprint(args...)
	log(DebugLevel, msg)
}

func getLogColor(l logLevel) (string, string) {
	switch l {
	case DebugLevel:
		return colors.DebugTextColor, "DEBUG"
	case InfoLevel:
		return colors.InfoTextColor, "INFO"
	case WarningLevel:
		return colors.WarningTextColor, "WARNING"
	case ErrorLevel:
		return colors.ErrorTextColor, "ERROR"
	case SuccessLevel:
		return colors.SuccessTextColor, "SUCCESS"
	default:
		return colors.WhiteTextColor, "LOG"
	}
}
