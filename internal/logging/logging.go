// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/logging/logging.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains a logger for the csgo sync application.
*/

package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Level int

const (
	ErrLvl Level = iota
	InfoLvl
	DebugLvl
)

type Log struct {
	lvl   Level
	err   *log.Logger
	info  *log.Logger
	debug *log.Logger
}

func Init(lvl Level, logFile io.WriteCloser) *Log {
	if lvl == ErrLvl { // level has to be at least 1 or "info"
		lvl = InfoLvl
	}
	return &Log{
		lvl:   lvl,
		err:   log.New(logFile, "[ERROR]: ", log.LstdFlags),
		info:  log.New(logFile, "[INFO]:  ", log.LstdFlags),
		debug: log.New(logFile, "[DEBUG]: ", log.LstdFlags),
	}
}

func (l *Log) Err(msg ...interface{}) {
	l.doLog(ErrLvl, msg...)
}

func (l *Log) Info(msg ...interface{}) {
	l.doLog(InfoLvl, msg...)
}

func (l *Log) Debug(msg ...interface{}) {
	l.doLog(DebugLvl, msg...)
}

func (l *Log) doLog(lvl Level, msg ...interface{}) {
	if lvl <= l.lvl {
		switch lvl {
		case ErrLvl:
			l.err.Println(msg...)
		case InfoLvl:
			l.info.Println(msg...)
		case DebugLvl:
			l.debug.Println(msg...)
		default:
			fmt.Println("no message!")
		}
	}
}

func OpenLogFile(fileName string) (io.WriteCloser, error) {
	switch strings.ToLower(fileName) {
	case "stderr", "":
		return os.Stderr, nil
	case "stdout":
		return os.Stdout, nil
	default:
		return os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	}
}
