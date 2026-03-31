package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Level represents logging level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR"}

// Logger is a thread-safe logger with level support
type Logger struct {
	mu     sync.Mutex
	level  Level
	output *os.File
}

// Global logger instance
var defaultLogger *Logger
var once sync.Once

// Init initializes the global logger
func Init(debug bool) {
	once.Do(func() {
		level := INFO
		if debug {
			level = DEBUG
		}
		defaultLogger = &Logger{
			level:  level,
			output: os.Stdout,
		}
	})
}

// SetOutput sets the output file for logging
func SetOutput(f *os.File) {
	if defaultLogger != nil {
		defaultLogger.output = f
	}
}

// SetLevel sets the logging level
func SetLevel(level Level) {
	if defaultLogger != nil {
		defaultLogger.level = level
	}
}

// formatHeader creates a log header with timestamp and level
func formatHeader(level Level, prefix string) string {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	return fmt.Sprintf("[%s] [%s] %s ", now, levelNames[level], prefix)
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	if defaultLogger != nil && defaultLogger.level <= DEBUG {
		defaultLogger.log(DEBUG, fmt.Sprintf(format, args...))
	}
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	if defaultLogger != nil && defaultLogger.level <= INFO {
		defaultLogger.log(INFO, fmt.Sprintf(format, args...))
	}
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	if defaultLogger != nil && defaultLogger.level <= WARN {
		defaultLogger.log(WARN, fmt.Sprintf(format, args...))
	}
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	if defaultLogger != nil && defaultLogger.level <= ERROR {
		defaultLogger.log(ERROR, fmt.Sprintf(format, args...))
	}
}

// log writes the log entry
func (l *Logger) log(level Level, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	header := formatHeader(level, "")
	if level == DEBUG {
		header = fmt.Sprintf("\033[36m%s\033[0m", header) // cyan for debug
	} else if level == ERROR {
		header = fmt.Sprintf("\033[31m%s\033[0m", header) // red for error
	} else if level == WARN {
		header = fmt.Sprintf("\033[33m%s\033[0m", header) // yellow for warn
	}
	fmt.Fprintf(l.output, "%s%s\n", header, msg)
}

// WithPrefix returns a function that logs with a specific prefix
func WithPrefix(prefix string) func(format string, args ...interface{}) {
	return func(format string, args ...interface{}) {
		Info(prefix+" "+format, args...)
	}
}
