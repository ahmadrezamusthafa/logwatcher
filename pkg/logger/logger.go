package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	syslog "github.com/RackSec/srslog"
	"golang.org/x/net/context"
)

const (
	LogLevelTrace   = "TRACE"
	LogLevelDebug   = "DEBUG"
	LogLevelInfo    = "INFO"
	LogLevelWarning = "WARNING"
	LogLevelError   = "ERROR"
	LogLevelPanic   = "PANIC"
)

var logLevel = map[string]int{
	LogLevelTrace:   5,
	LogLevelDebug:   4,
	LogLevelInfo:    3,
	LogLevelWarning: 2,
	LogLevelError:   1,
	LogLevelPanic:   0,
}

var ptSystemName string

var activeLogLevel = strings.ToUpper(os.Getenv("LOG_LEVEL"))

func parseLogLevel() string {
	switch activeLogLevel {
	case LogLevelTrace, LogLevelDebug, LogLevelInfo, LogLevelWarning, LogLevelError, LogLevelPanic:
	default:
		activeLogLevel = LogLevelTrace
	}
	return activeLogLevel
}

func getActiveLogLevel() int {
	return logLevel[activeLogLevel]
}

func extractReqID(ctx context.Context) string {
	requestIDKey := "x-request-id"
	str, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return str
}

// SetupLogger creates logger instance to log to PaperTrail and Console. Should only called once in main function.
func SetupLogger(ptHost string, ptPort string) {
	activeLogLevel = parseLogLevel()
	log.SetPrefix("")
	log.SetFlags(0)

	if ptHost != "" {
		hostname, _ := os.Hostname()
		ptEndpoint := fmt.Sprintf("%s:%s", ptHost, ptPort)
		ptWriter, err := syslog.Dial("udp", ptEndpoint, syslog.LOG_INFO, hostname)

		if err != nil {
			log.Fatal("Can't connect to PaperTrail ...")
		}

		log.SetOutput(io.MultiWriter(os.Stdout, ptWriter))
	} else {
		log.Print("No papertrail transport detected. Logger only use local stdout")
	}
}

func formatterRFC3164(p syslog.Priority, hostname, tag, content string) string {
	return syslog.RFC3164Formatter(p, ptSystemName, hostname, content)
}

// SetupLoggerAuto creates logger instance to log to PaperTrail automatically without specifying PT HOST and PORT
func SetupLoggerAuto(appName string, ptEndpoint string) {
	activeLogLevel = parseLogLevel()
	log.SetPrefix("")
	log.SetFlags(0)

	if appName != "" && ptEndpoint != "" {
		ptSystemName = appName
		hostname, _ := os.Hostname()

		ptWriter, err := syslog.Dial("udp", ptEndpoint, syslog.LOG_INFO, hostname)

		if err != nil {
			log.Fatalf("Can't connect to PaperTrail: %s", err.Error())
		}

		ptWriter.SetFormatter(formatterRFC3164)

		log.SetOutput(io.MultiWriter(os.Stdout, ptWriter))
	} else {
		log.Print("Logger configured to use only local stdout")
	}
}

// Warn prints warning message to logs
func Warn(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelWarning] {
		message := fmt.Sprintf("WARN: "+format, v...)
		log.Print(message)
	}
}

// Trace prints trace message to logs
func Trace(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelTrace] {
		message := fmt.Sprintf("TRACE: "+format, v...)
		log.Print(message)
	}
}

// Debug prints debug message to logs
func Debug(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelDebug] {
		message := fmt.Sprintf("DEBUG: "+format, v...)
		log.Print(message)
	}
}

// Info prints info message to logs
func Info(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelInfo] {
		message := fmt.Sprintf("INFO: "+format, v...)
		log.Print(message)
	}
}

// Err prints error message to logs without stacktrace
func Err(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelError] {
		message := []interface{}{fmt.Sprintf("ERROR: "+format, v...)}
		log.Print(message...)
	}
}

// Panic prints panic message to logs
func Panic(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelPanic] {
		message := []interface{}{fmt.Sprintf("PANIC: "+format, v...)}
		log.Print(message...)
	}
}

// Fatal calls Err and then os.Exit(1)
func Fatal(format string, v ...interface{}) {
	Err(format, v...)
	os.Exit(1)
}

// WarnContext prints warning message to logs
func WarnContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelWarning] {
		message := fmt.Sprintf("WARN: "+extractReqID(ctx)+" - "+format, v...)
		log.Print(message)
	}
}

// TraceContext prints trace message to logs
func TraceContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelTrace] {
		message := fmt.Sprintf("TRACE: ReqID "+extractReqID(ctx)+" - "+format, v...)
		log.Print(message)
	}
}

// DebugContext prints debug message to logs
func DebugContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelDebug] {
		message := fmt.Sprintf("DEBUG: ReqID "+extractReqID(ctx)+" - "+format, v...)
		log.Print(message)
	}
}

// InfoContext prints info message to logs
func InfoContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelInfo] {
		message := fmt.Sprintf("INFO: "+extractReqID(ctx)+" - "+format, v...)
		log.Print(message)
	}
}

// ErrContext prints error message to logs without stacktrace
func ErrContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelError] {
		message := []interface{}{fmt.Sprintf("ERROR: "+extractReqID(ctx)+" - "+format, v...)}
		log.Print(message...)
	}
}

// PanicContext prints panic message to logs
func PanicContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelPanic] {
		message := []interface{}{fmt.Sprintf("PANIC: ReqID "+extractReqID(ctx)+" - "+format, v...)}
		log.Print(message...)
	}
}

// FatalContext calls Err and then os.Exit(1)
func FatalContext(ctx context.Context, format string, v ...interface{}) {
	ErrContext(ctx, format, v...)
	os.Exit(1)
}
