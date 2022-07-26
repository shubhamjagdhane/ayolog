package mylog

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

type Entry struct {
	Logger *Logger

	Data Fields

	Time time.Time

	Level Level

	//	Caller *runtime.Frame

	Message string

	Buffer *bytes.Buffer
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		Data:   make(Fields, 6),
	}
}

func (entry *Entry) write() {
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}

	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	if _, err := entry.Logger.Out.Write(serialized); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}

func (entry *Entry) Dup() *Entry {
	data := make(Fields, len(entry.Data))
	for k, v := range entry.Data {
		data[k] = v
	}

	return &Entry{
		Logger: entry.Logger, Data: data, Time: entry.Time,
	}
}

func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	}
	return "something"
}

func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case DebugLevel:
		return []byte("debug"), nil
	case InfoLevel:
		return []byte("info"), nil
	case ErrorLevel:
		return []byte("error"), nil
	case FatalLevel:
		return []byte("fatal"), nil
	case PanicLevel:
		return []byte("panic"), nil
	}

	return nil, fmt.Errorf("Not a valid Level: %v", level)
}

func (entry *Entry) Log(level Level, args ...interface{}) {
	if entry.Logger.IsLevelEnable(level) {
		entry.log(level, fmt.Sprint(args...))
	}
}

func (entry *Entry) log(level Level, msg string) {
	var buffer *bytes.Buffer

	newEntry := entry.Dup()

	if newEntry.Time.IsZero() {
		newEntry.Time = time.Now()
	}

	newEntry.Level = level
	newEntry.Message = msg

	buffer = getBuffer()
	defer func() {
		newEntry.Buffer = nil
		putBuffer(buffer)
	}()

	buffer.Reset()
	newEntry.Buffer = buffer

	newEntry.write()
	newEntry.Buffer = nil

	if level <= PanicLevel {
		panic(newEntry)
	}
}
