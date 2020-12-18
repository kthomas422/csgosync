// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/logging/logging.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains a logger for the csgo sync application.
*/

package csgolog

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kthomas422/csgosync/config"

	logger "github.com/kthomas422/json-logger"
)

type CsgoLogger struct {
	file io.WriteCloser
	*logger.Logger
}

func InitLogger(logPath string) (*CsgoLogger, error) {
	switch logPath {
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

func (cl CsgoLogger) Close() error {
	if cl.file == os.Stderr || cl.file == os.Stdout {
		return nil
	}
	return cl.file.Close()
}

func (cl CsgoLogger) WebRequest(request *http.Request) error {
	return cl.Info(logger.ServerRequestLog{
		Origin: request.Host,
		URI:    request.RequestURI,
		Header: request.Header,
	})
}

func (cl CsgoLogger) Config(config config.ServerConfig) error {
	return cl.Info(logger.ConfigLog{Config: config})
}

func (cl CsgoLogger) Simple(msg string) {
	if err := cl.Info(logger.SimpleLog(msg)); err != nil {
		panic(fmt.Errorf("failed to write to logger: %w", err))
	}
}

func (cl CsgoLogger) Err(msg string, err error) {
	if err := cl.Error(logger.Simple{
		Msg: fmt.Sprint(msg, err),
	}); err != nil {
		panic(fmt.Errorf("failed to write to logger: %w", err))
	}
}
