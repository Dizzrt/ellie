package log

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var DefaultLogger = NewStdLogger(log.Writer())

var global = &loggerAppliance{}

type loggerAppliance struct {
	LogWriter
	lock sync.RWMutex
}

func init() {
	global.SetLogger(DefaultLogger)
}

func (a *loggerAppliance) SetLogger(writer LogWriter) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.LogWriter = writer
}

func SetLogger(writer LogWriter) {
	global.SetLogger(writer)
}

func GetLogger() LogWriter {
	global.lock.RLock()
	defer global.lock.RUnlock()
	return global.LogWriter
}

func Debug(a ...any) {
	global.Write(LevelDebug, DefaultMessageKey, fmt.Sprint(a...))
}

func Debugf(format string, a ...any) {
	global.Write(LevelDebug, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Debugw(keyvals ...any) {
	global.Write(LevelDebug, keyvals...)
}

func Info(a ...any) {
	global.Write(LevelInfo, DefaultMessageKey, fmt.Sprint(a...))
}

func Infof(format string, a ...any) {
	global.Write(LevelInfo, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Infow(keyvals ...any) {
	global.Write(LevelInfo, keyvals...)
}

func Warn(a ...any) {
	global.Write(LevelWarn, DefaultMessageKey, fmt.Sprint(a...))
}

func Warnf(format string, a ...any) {
	global.Write(LevelWarn, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Warnw(keyvals ...any) {
	global.Write(LevelWarn, keyvals...)
}

func Error(a ...any) {
	global.Write(LevelError, DefaultMessageKey, fmt.Sprint(a...))
}

func Errorf(format string, a ...any) {
	global.Write(LevelError, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Errorw(keyvals ...any) {
	global.Write(LevelError, keyvals...)
}

func Fatal(a ...any) {
	global.Write(LevelFatal, DefaultMessageKey, fmt.Sprint(a...))
	os.Exit(1)
}

func Fatalf(format string, a ...any) {
	global.Write(LevelFatal, DefaultMessageKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

func Fatalw(keyvals ...any) {
	global.Write(LevelFatal, keyvals...)
	os.Exit(1)
}
