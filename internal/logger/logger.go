package logger

import (
	"io"
	"log/slog"
	"os"
)

type Config struct {
	accessLogFilename string `yaml:"accessLogFilename"`
	errorLogFilename  string `yaml:"errorLogFilename"`
}

type Logger struct {
	conf         Config
	accessLogger *slog.Logger
	errorLogger  *slog.Logger
}

func New(conf Config) (*Logger, error) {

	accessLogFile, err := os.OpenFile(conf.accessLogFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	errorLogFile, err := os.Create(conf.errorLogFilename)
	if err != nil {
		return nil, err
	}

	mw := io.MultiWriter(errorLogFile, os.Stdout)

	accesLoggerHandler := slog.NewTextHandler(accessLogFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	errorLoggerHandler := slog.NewTextHandler(mw, &slog.HandlerOptions{
		Level: slog.LevelError,
	})

	return &Logger{
		conf:         conf,
		accessLogger: slog.New(accesLoggerHandler),
		errorLogger:  slog.New(errorLoggerHandler),
	}, nil
}

func (l *Logger) Error(err error) {
	l.errorLogger.Error(err.Error())
}

func (l *Logger) Info(msg string) {
	l.accessLogger.Info(msg)
}
