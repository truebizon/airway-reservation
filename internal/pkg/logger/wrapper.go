package logger

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Wrapper struct {
	Fields []zap.Field
}

func checkLevel(levelStr string) bool {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		return false
	}
	if logLevel > level {
		return false
	}
	return true
}

// NewWrapper can additional fields.
// ex. Use this when you want to add a common value in the scope of a context, such as an API request.
func NewWrapper(fields ...zap.Field) *Wrapper {
	return &Wrapper{Fields: fields}
}

func (w *Wrapper) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper().Debug(msg, fields...)
}

func (w *Wrapper) Info(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper().Info(msg, fields...)
}

func (w *Wrapper) Warn(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper().Warn(msg, fields...)
}

func (w *Wrapper) Error(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper().Error(msg, fields...)
}

func (w *Wrapper) Fatal(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper().Fatal(msg, fields...)
}

func (w *Wrapper) Debugf(msg string, fields ...interface{}) {
	wrapper().Debug(fmt.Sprintf(msg, fields...), w.Fields...)
}

func (w *Wrapper) Infof(msg string, fields ...interface{}) {
	wrapper().Info(fmt.Sprintf(msg, fields...), w.Fields...)
}

func (w *Wrapper) Warnf(msg string, fields ...interface{}) {
	wrapper().Warn(fmt.Sprintf(msg, fields...), w.Fields...)
}

func (w *Wrapper) Errorf(msg string, err error, fields ...interface{}) {
	wrapper().Error(fmt.Sprintf(msg, fields...), addErrorField(w.Fields, err)...)
}

func (w *Wrapper) Fatalf(msg string, err error, fields ...interface{}) {
	wrapper().Fatal(fmt.Sprintf(msg, fields...), addErrorField(w.Fields, err)...)
}

// Sync wrapper of Zap's Sync.
func Sync() {
	Info("FLUSH_LOG_BUFFER")
	if err := zapLogger.Sync(); err != nil {
		log.Fatal(err)
	}
}

// SyncWhenStop flush log buffer. when interrupt or terminated.
func SyncWhenStop() {
	c := make(chan os.Signal, 1)

	go func() {
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		s := <-c

		sigCode := 0
		switch s.String() {
		case "interrupt":
			sigCode = 2
		case "terminated":
			sigCode = 15
		}

		Info(fmt.Sprintf("GOT_SIGNAL_%v", strings.ToUpper(s.String())))
		Sync() // flush log buffer
		os.Exit(128 + sigCode)
	}()
}

// Debug is Wrapper of Zap's Debug.
// Outputs a short log to the console. Detailed json log output to log file.
func Debug(msg string, fields ...zap.Field) {
	wrapper().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	wrapper().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	wrapper().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	wrapper().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	wrapper().Fatal(msg, fields...)
}

// Debugf is Outputs a Debug log with formatted error.
func Debugf(msg string, fields ...interface{}) {
	wrapper().Debug(fmt.Sprintf(msg, fields...))
}

func Infof(msg string, fields ...interface{}) {
	wrapper().Info(fmt.Sprintf(msg, fields...))
}

func Warnf(msg string, fields ...interface{}) {
	wrapper().Warn(fmt.Sprintf(msg, fields...))
}

func Errorf(msg string, err error, fields ...interface{}) {
	wrapper().Error(fmt.Sprintf(msg, fields...), addErrorField([]zap.Field{}, err)...)
}

func Fatalf(msg string, err error, fields ...interface{}) {
	wrapper().Fatal(fmt.Sprintf(msg, fields...), addErrorField([]zap.Field{}, err)...)
}

func wrapper() *zap.Logger {
	checkInit()
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func addErrorField(fields []zap.Field, err error) []zap.Field {
	return append(fields, zap.String("error", fmt.Sprintf("%+v", err)))
}

func checkInit() {
	if zapLogger == nil {
		log.Fatal("The logger is not initialized. InitLogger() must be called.")
	}
}

// Dump is Wrapper of spew.Dump()
func Dump(i interface{}) {
	if !checkLevel("DEBUG") {
		return
	}
	spew.Dump(i)
}
