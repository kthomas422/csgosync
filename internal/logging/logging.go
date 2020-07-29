package logging

import (
	"io"
	"log"
)

type Level int

const (
	ErrLvl Level = iota
	InfoLvl
	DebugLvl
)

type Log struct {
	lvl     Level
	logFile io.Writer
	err     *log.Logger
	info    *log.Logger
	debug   *log.Logger
}

func Init(lvl Level, logFile io.Writer) *Log {
	return &Log{
		lvl:     lvl,
		logFile: logFile,
		err:     log.New(logFile, "[ERROR]: ", log.LstdFlags),
		info:    log.New(logFile, "[INFO]:  ", log.LstdFlags),
		debug:   log.New(logFile, "[DEBUG]: ", log.LstdFlags),
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
		}
	}
}
