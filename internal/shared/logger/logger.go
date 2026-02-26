package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogNum int

const (
	Debug LogNum = iota
	Info
	Warning
	Error
)
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

type Logger struct {
	MinLevel LogNum
	termLog  *log.Logger
	fileLog  *log.Logger
}

func NewLogger(minLevel LogNum, logDir string, maxDays int) *Logger {
	os.MkdirAll(logDir, 0755)

	today := time.Now().Format("2006-01-02")
	logFile := logDir + "chat-app-" + today + ".log"

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	flags := log.Ldate | log.Ltime | log.Lshortfile
	return &Logger{
		MinLevel: minLevel,
		termLog:  log.New(os.Stdout, "", flags),
		fileLog:  log.New(file, "", flags),
	}
}

func (l *Logger) writeLog(level LogNum, colorPrefix, plainPrefix, format string, v ...any) {

	if l.MinLevel <= level {

		msg := fmt.Sprintf(format, v...)

		l.termLog.Output(3, colorPrefix+msg+colorReset)
		l.fileLog.Output(3, plainPrefix+msg+colorReset)

	}

}

func (l *Logger) Debug(format string, v ...any) {
	l.writeLog(Debug, colorCyan+"[DEBUG] ", "[DEBUG] ", format, v...)
}

func (l *Logger) Info(format string, v ...any) {
	l.writeLog(Info, colorBlue+"[INFO] ", "[INFO] ", format, v...)
}

func (l *Logger) Warning(format string, v ...any) {
	l.writeLog(Warning, colorYellow+"[WARN] ", "[WARN] ", format, v...)
}

func (l *Logger) Error(format string, v ...any) {
	l.writeLog(Error, colorRed+"[ERROR] ", "[ERROR] ", format, v...)
}
