package log

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	httpUtils "gitlab.com/docshade/common/http"
	"gitlab.com/docshade/common/utils"

	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log       Logger
	logConfig LoggerConfig
)

const (
	//Debug has verbose message
	Debug = "debug"
	//Info is default log level
	Info = "info"
	//Warn is for logging messages about possible issues
	Warn = "warn"
	//Error is for logging errors
	Error = "error"
	//Fatal is for logging fatal messages. The sytem shutsdown after logging the message.
	Fatal = "fatal"
	// ContextRequestIDKey - Context Key to store RequestID
	ContextRequestIDKey requestID = "request_id"
)

type requestID string

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Debugc(ctx context.Context, format string, args ...interface{})
	Infoc(ctx context.Context, format string, args ...interface{})
	Warnc(ctx context.Context, format string, args ...interface{})
	Errorc(ctx context.Context, format string, args ...interface{})
	Fatalc(ctx context.Context, format string, args ...interface{})
	Panicc(ctx context.Context, format string, args ...interface{})

	EchoLoggingMiddleware() echo.MiddlewareFunc
	RequestIDMiddleware(next http.Handler) http.Handler
	WithFields(keyValues Fields) Logger
}

// LoggerConfig stores the config for the logger
type LoggerConfig struct {
	LogSkipURNs    []string `mapstructure:"skip_urns"`
	LogLevel       string   `mapstructure:"level"`
	LogMemoryUsage bool     `mapstructure:"memory_usage"`
	LogJSON        bool     `mapstructure:"json"`
	LogShowCaller  bool     `mapstructure:"show_caller"`
	LogCallerSkip  int      `mapstructure:"caller_skip"`
}

// Bytes to MiB
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// InitLog creates new logger
func InitLog(cfg LoggerConfig) {
	logConfig = cfg

	rand.Seed(time.Now().UTC().UnixNano())

	if logConfig.LogLevel == "" {
		logConfig.LogLevel = "info"
	}

	NewLogger(logConfig)
}

// NewLogger -
func NewLogger(cfg LoggerConfig) {
	var cores []zapcore.Core
	level := getZapLevel(cfg.LogLevel)
	writer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(getEncoder(cfg.LogJSON), writer, level)
	cores = append(cores, core)
	combinedCore := zapcore.NewTee(cores...)

	var logger *zap.SugaredLogger
	if logConfig.LogShowCaller {
		logger = zap.New(combinedCore,
			zap.AddCallerSkip(cfg.LogCallerSkip),
			zap.AddCaller(),
		).Sugar()
	} else {
		logger = zap.New(combinedCore).Sugar()
	}

	log = &zapLogger{
		sugaredLogger: logger,
	}

	if logConfig.LogMemoryUsage {
		log.Debugf("Memory Usage legend: A = Alloc, TA = TotalAlloc, S = System, NGC = NextGCHeap, GC = GCCycles")
	}
}

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
	C             LoggerConfig
}

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = fixedLengthCallerEncoder
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// fixedLengthCallerEncoder ensures that caller takes exactly 30 chars to align the columns in the output
func fixedLengthCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if logConfig.LogJSON {
		enc.AppendString(caller.TrimmedPath())
	} else {
		enc.AppendString(fmt.Sprintf("%-30s", caller.TrimmedPath()))
	}
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Debug:
		return zapcore.DebugLevel
	case Error:
		return zapcore.ErrorLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// memLog - return memory usage string
func memLog() string {
	if logConfig.LogMemoryUsage {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		// For info on each, see: https://golang.org/pkg/runtime/#MemStats
		return fmt.Sprintf(" /%v/%v/%v/%v/%v", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), bToMb(m.NextGC), m.NumGC)
	}
	return ""
}

/* Default log functions
========================================================================= */

// Warne is for logging the concrete error
func Warne(err error) {
	log.Warnf("%s", err)
}

// Errorf is for logging errors
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf is for logging fatal messages. The sytem shutsdown after logging the message.
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Panicf is for logging panic messages. The sytem will panics after logging the message.

// WithFields creates a child logger and adds structured context to it. Fields added to the child don't affect the parent.
func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

// RequestIDMiddleware -
func RequestIDMiddleware(next http.Handler) http.Handler {
	return log.RequestIDMiddleware(next)
}

func EchoLogger() echo.MiddlewareFunc {
	return log.EchoLoggingMiddleware()
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format+memLog(), args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	if strings.Contains(format, "\n") {
		frmt := fmt.Sprintf(format, args...)
		for _, f := range strings.Split(frmt, "\n") {
			l.sugaredLogger.Infof(f + memLog())
		}
	} else {
		l.sugaredLogger.Infof(format+memLog(), args...)
	}
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format+memLog(), args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format+memLog(), args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format+memLog(), args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format+memLog(), args...)
}

func (l *zapLogger) Debugc(ctx context.Context, format string, args ...interface{}) {
	l.sugaredLogger.Debugf(GetRequestID(ctx)+format+memLog(), args...)
}

func (l *zapLogger) Infoc(ctx context.Context, format string, args ...interface{}) {
	l.sugaredLogger.Infof(GetRequestID(ctx)+format+memLog(), args...)
}

func (l *zapLogger) Warnc(ctx context.Context, format string, args ...interface{}) {
	l.sugaredLogger.Warnf(GetRequestID(ctx)+format+memLog(), args...)
}

func (l *zapLogger) Errorc(ctx context.Context, format string, args ...interface{}) {
	l.sugaredLogger.Errorf(GetRequestID(ctx)+format+memLog(), args...)
}

func (l *zapLogger) Fatalc(ctx context.Context, format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(GetRequestID(ctx)+format+memLog(), args...)
}

func (l *zapLogger) Panicc(ctx context.Context, format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(GetRequestID(ctx)+format+memLog(), args...)
}

/* Tag related log functions
========================================================================= */

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{sugaredLogger: newLogger}
}

// RequestIDGen - generate random request ID (from "1000" to "ZZZZ")
func RequestIDGen() string {
	return strings.ToUpper(strconv.FormatInt(rand.Int63n(1632959)+46656, 36))
}

// GetRequestID - return request ID from Context
func GetRequestID(ctx context.Context) string {
	if ctx != nil {
		if ctxRqID, ok := ctx.Value(ContextRequestIDKey).(string); ok {
			return "[" + ctxRqID + "] "
		}
	}
	return ""
}

// Middleware functions
// RequestIDMiddleware - injects random request id in http.Request Context
func (l *zapLogger) RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := RequestIDGen()
		ctx := context.WithValue(r.Context(), ContextRequestIDKey, requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// responseWriter wrapper for http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	wroteheader bool
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteheader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteheader = true
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size = size
	if rw.status == 0 {
		rw.status = 200
	}
	return size, err
}

func containsString(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// LoggingMiddleware - to log incoming HTTP requests
func (l *zapLogger) EchoLoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			r := c.Request()
			w := c.Response()

			t0 := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}

			if containsString(r.URL.EscapedPath(), logConfig.LogSkipURNs) {
				return
			}
			var q string
			if len(r.URL.Query()) > 0 {
				q = "?" + r.URL.RawQuery
			}

			refHost := "-"
			referer := r.Referer()
			if referer != "" {
				refURL, err := url.Parse(referer)
				if err == nil && refURL.Hostname() != "" {
					refHost = refURL.Scheme + "://" + strings.TrimSuffix(refURL.Hostname(), "/")
				}
			} else {
				referer = "-"
			}

			logger := l.WithFields(Fields{
				"requestId": GetRequestID(r.Context()),
				"method":    r.Method,
				"headers":   getHeaders(c),
				"url":       r.URL.EscapedPath() + q,
				"status":    w.Status,
				"latency":   time.Since(t0),
				"proto":     r.Proto,
				"refHost":   refHost,
				"referer":   referer,
				"size":      w.Size,
				"userAgent": r.UserAgent(),
				"memLog":    memLog(),
			})
			logger.Infoc(r.Context(), "Served")
			return nil
		}
	}
}

func getHeaders(ctx echo.Context) http.Header {
	rh := ctx.Request().Header

	for k, v := range rh {
		if _, ok := httpUtils.MaskHeaders[k]; ok {
			masked := make([]string, 0, len(v))
			for _, s := range v {
				masked = append(masked, utils.MaskText(s))
			}
			rh[k] = masked
		}
	}

	return rh
}
