package logruslog

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/common/vipercfg"
)

const RootLoggerName = "ROOT"
const DefaultLoggerLevel = "INFO"

// Mapping between logger name and level
type LoggerLevelMapping map[string]string

type LoggerFactory struct {
	LoggerLevels LoggerLevelMapping
}

// Gets logger with pre-configured level
func (factory *LoggerFactory) GetLogger(name string) *log.Logger {
	level, ok := factory.LoggerLevels[name]
	if ok {
		return NewDefaultLogger(level)
	}

	rootLevel, rootOk := factory.LoggerLevels[RootLoggerName]
	if rootOk {
		return NewDefaultLogger(rootLevel)
	}

	return NewDefaultLogger(DefaultLoggerLevel)
}

// Initialize a logger with string value of level
func NewDefaultLogger(logLevelValue string) *log.Logger {
	newLogger := log.New()
	newLogger.Formatter = newDefulatFormatter()
	newLogger.Level = logLevel(logLevelValue)

	return newLogger
}

// Sets log level by various string
//
// "debug", "Debug", "DEBUG"
// "info", "Info", "INFO"
// "warn", "Warn", "WARN"
// "error", "Error", "ERROR"
// "fatal", "Fatal", "FATAL"
// "panic", "Panic", "PANIC"
//
// Otherwise, use **INFO** level
func SetLogLevelByString(logLevelValue string) {
	log.SetFormatter(newDefulatFormatter())
	log.SetLevel(logLevel(logLevelValue))
}

func newDefulatFormatter() *log.TextFormatter {
	return &log.TextFormatter{FullTimestamp: true}
}

type funcInfo struct {
	file string
	line int
	name string
}

func (f funcInfo) String() string {
	return fmt.Sprintf("%s (%s:%d)", path.Base(f.name), path.Base(f.file), f.line)
}

type stackHook struct{}

func (hook stackHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook stackHook) Fire(entry *log.Entry) error {
	var skipFrames int
	if len(entry.Data) == 0 {
		// When WithField(s) is not used, we have 8 logrus frames to skip.
		skipFrames = 8
	} else {
		// When WithField(s) is used, we have 6 logrus frames to skip.
		skipFrames = 6
	}
	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(skipFrames, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/Sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			f := funcInfo{
				file: file,
				line: line,
				name: name,
			}
			entry.Data["func"] = f
			break
		}
	}
	return nil
}

func logLevel(l string) log.Level {
	switch strings.ToLower(l) {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

func Init() {
	SetLogLevelByString(vipercfg.Config().GetString("logLevel"))
	log.AddHook(stackHook{})
}
