package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	AccessLogFilename string `yaml:"accessLogFilename"`
	ErrorLogFilename  string `yaml:"errorLogFilename"`
}

type Logger struct {
	conf         Config
	accessLogger *slog.Logger
	errorLogger  *slog.Logger
}

func New(conf Config) (*Logger, error) {

	fmt.Printf("%+v", conf)

	accessLogFile := &lumberjack.Logger{
		Filename:   conf.AccessLogFilename,
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}

	errorLogFile := &lumberjack.Logger{
		Filename:   conf.ErrorLogFilename,
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}

	mw := io.MultiWriter(errorLogFile, os.Stdout)

	accessLoggerHandler := slog.NewTextHandler(accessLogFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	errorLoggerHandler := slog.NewTextHandler(mw, &slog.HandlerOptions{
		Level: slog.LevelError,
	})

	return &Logger{
		conf:         conf,
		accessLogger: slog.New(accessLoggerHandler),
		errorLogger:  slog.New(errorLoggerHandler),
	}, nil
}

func (l *Logger) Error(msg ...string) {
	l.errorLogger.Error(strings.Join(msg, " "))
}

func (l *Logger) Info(msg ...string) {
	l.accessLogger.Info(strings.Join(msg, " "))
}
