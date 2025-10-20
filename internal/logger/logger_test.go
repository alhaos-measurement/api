package logger

import "testing"

func TestLogger_Info(t *testing.T) {

	cfg := Config{
		accessLogFilename: "access.log",
		errorLogFilename:  "error.log",
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("hello world")
}
