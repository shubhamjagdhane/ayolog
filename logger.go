package mylog

import (
	"io"
	"os"
	"sync"
	"sync/atomic"
)

type Level uint32

type Fields map[string]interface{}

type Logger struct {
	Out io.Writer

	Formatter Formatter

	//	ReportCaller bool

	Level Level

	mu MutexWrap

	entryPool sync.Pool
}

func New() *Logger {
	return &Logger{
		Out:       os.Stderr,
		Formatter: new(JSONFormatter),
		Level:     DebugLevel,
	}
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

func (logger *Logger) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{}
	logger.entryPool.Put(entry)
}

type MutexWrap struct {
	lock    sync.Mutex
	disable bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disable {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disable {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disable = true
}

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

func (logger *Logger) level() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.Level)))
}

func (logger *Logger) IsLevelEnable(level Level) bool {
	return logger.level() >= level
}

func (logger *Logger) Log(level Level, args ...interface{}) {
	if logger.IsLevelEnable(level) {
		entry := logger.newEntry()
		entry.Log(level, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Log(DebugLevel, args...)
}
