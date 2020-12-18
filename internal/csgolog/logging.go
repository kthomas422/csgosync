// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/logging/logging.go
	Language:	Go 1.15
	Dev Env:	Linux 5.9

	This file contains a logger for the csgo sync application.
*/

package csgolog

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kthomas422/csgosync/config"

	logger "github.com/kthomas422/json-logger"
)

// Wrapper for the logging file to close later and the logging package struct
type CsgoLogger struct {
	file io.WriteCloser
	*logger.Logger
}

// InitLogger determines if the log file is stderr/stdout or a file and opens it if a file
func InitLogger(logPath string) (*CsgoLogger, error) {
	switch strings.ToLower(logPath) {
	case "stderr":
		return &CsgoLogger{
			file:   os.Stderr,
			Logger: logger.NewLogger(os.Stderr),
		}, nil
	case "stdout":
		return &CsgoLogger{
			file:   os.Stdout,
			Logger: logger.NewLogger(os.Stdout),
		}, nil
	default:
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return nil, err
		}
		return &CsgoLogger{
			file:   file,
			Logger: logger.NewLogger(file),
		}, nil
	}
}

// Close closes the log file if its a file and not stdout or stderr
func (cl CsgoLogger) Close() error {
	if cl.file == os.Stderr || cl.file == os.Stdout {
		return nil
	}
	return cl.file.Close()
}

// WebRequest takes in a web request and makes a log entry for it
func (cl CsgoLogger) WebRequest(request *http.Request) {
	_ = cl.Info(logger.ServerRequestLog{
		Origin: request.Host,
		URI:    request.RequestURI,
		Header: request.Header,
	})
}

// Config takes in the server config and logs it
func (cl CsgoLogger) Config(config config.ServerConfig) error {
	return cl.Info(logger.ConfigLog{Config: config})
}

// Simple takes in a string "message" and logs it
func (cl CsgoLogger) Simple(msg string) {
	if err := cl.Info(logger.SimpleLog(msg)); err != nil {
		panic(fmt.Errorf("failed to write to logger: %w", err))
	}
}

// Err takes in a string message and an error and wraps them together to log
func (cl CsgoLogger) Err(msg string, err error) {
	if err := cl.Error(logger.Simple{
		Msg: fmt.Sprint(msg, err),
	}); err != nil {
		panic(fmt.Errorf("failed to write to logger: %w", err))
	}
}
