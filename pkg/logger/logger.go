package logger

import (
	"log"
	"os"
	"strings"
)

// LogLevel represents log levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger represents a structured logger
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

var (
	defaultLogger *Logger
)

// NewLogger creates a new logger with specified level
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// InitDefault initializes the default logger
func InitDefault(level LogLevel) {
	defaultLogger = NewLogger(level)
}

// Debug logs debug messages
func (l *Logger) Debug(msg string, fields ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Printf("[DEBUG] %s %v", msg, fields)
	}
}

// Info logs info messages
func (l *Logger) Info(msg string, fields ...interface{}) {
	if l.level <= INFO {
		l.logger.Printf("[INFO] %s %v", msg, fields)
	}
}

// Warn logs warning messages
func (l *Logger) Warn(msg string, fields ...interface{}) {
	if l.level <= WARN {
		l.logger.Printf("[WARN] %s %v", msg, fields)
	}
}

// Error logs error messages
func (l *Logger) Error(msg string, fields ...interface{}) {
	if l.level <= ERROR {
		l.logger.Printf("[ERROR] %s %v", msg, fields)
	}
}

// Debugf logs formatted debug messages
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Printf("[DEBUG] "+format, args...)
	}
}

// Infof logs formatted info messages
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.level <= INFO {
		l.logger.Printf("[INFO] "+format, args...)
	}
}

// Warnf logs formatted warning messages
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.level <= WARN {
		l.logger.Printf("[WARN] "+format, args...)
	}
}

// Errorf logs formatted error messages
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.logger.Printf("[ERROR] "+format, args...)
	}
}

// Default logger functions
func Debug(msg string, fields ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, fields...)
	}
}

func Debugf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Errorf(format, args...)
	}
}

// ParseLogLevel parses string to LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}
