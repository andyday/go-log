package log

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	logger    = logrus.New()
	ctxFields []interface{}
)

type Formatter int

const (
	SimpleFormatter Formatter = iota
	TextFormatter
	JSONFormatter
)

type Level = logrus.Level

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = logrus.PanicLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = logrus.FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = logrus.ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = logrus.WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = logrus.InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = logrus.DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel = logrus.TraceLevel
)

type simpleFormatter struct{}

func (s *simpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if len(entry.Data) > 0 {
		data := strings.Builder{}
		for k, v := range entry.Data {
			data.WriteString(" | ")
			data.WriteString(k)
			data.WriteRune('=')
			sv, ok := v.(string)
			if !ok {
				sv = jsonString(v)
			}

			data.WriteString(sv)
		}
		return []byte(entry.Message + "  " + data.String() + "\n"), nil
	}
	return []byte(entry.Message + "\n"), nil
}

var formatMap = map[string]Formatter{
	"simple": SimpleFormatter,
	"text":   TextFormatter,
	"json":   JSONFormatter,
}

func FormatterFromName(name string) (f Formatter) {
	var ok bool
	name = strings.ToLower(name)
	if f, ok = formatMap[name]; ok {
		return f
	}
	return JSONFormatter
}

func Init(formatter Formatter, level Level, contextFields ...interface{}) {
	switch formatter {
	case JSONFormatter:
		logger.SetFormatter(new(logrus.JSONFormatter))
	case TextFormatter:
		logger.SetFormatter(new(logrus.TextFormatter))
	case SimpleFormatter:
		logger.SetFormatter(new(simpleFormatter))
	}
	logger.SetLevel(level)
	ctxFields = contextFields
}

func withContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}
	for _, f := range ctxFields {
		val := ctx.Value(f)
		if val != nil {
			fields[fmt.Sprintf("%v", f)] = val.(string)
		}
	}
	return logger.WithFields(fields)
}

type Fld interface {
	apply(fields logrus.Fields)
}

type fld struct {
	key   string
	value interface{}
}

func (f *fld) apply(fields logrus.Fields) {
	fields[f.key] = f.value
}

func Field(key string, value interface{}) Fld {
	if err, ok := value.(error); ok {
		value = err.Error()
	}
	return &fld{key: key, value: value}
}

func withFields(entry *logrus.Entry, flds []Fld) *logrus.Entry {
	fields := make(logrus.Fields)
	for _, f := range flds {
		f.apply(fields)
	}
	return entry.WithFields(fields)
}

// Info prints logs while attempting to JSON dump any non-primitive argument.
func Info(ctx context.Context, i interface{}, flds ...Fld) {
	withFields(withContext(ctx), flds).Info(i)
}

// Infof prints formatted logs while attempting to JSON dump any non-primitive argument.
func Infof(ctx context.Context, format string, a ...interface{}) {
	withContext(ctx).Infof(format, normalizeArgs(a)...)
}

// Warn prints logs while attempting to JSON dump any non-primitive argument.
func Warn(ctx context.Context, w interface{}, flds ...Fld) {
	withFields(withContext(ctx), flds).Warn(w)
}

// Warnf prints formatted logs while attempting to JSON dump any non-primitive argument.
func Warnf(ctx context.Context, format string, a ...interface{}) {
	withContext(ctx).Warnf(format, normalizeArgs(a)...)
}

// Error prints logs while attempting to JSON dump any non-primitive argument.
func Error(ctx context.Context, e interface{}, flds ...Fld) {
	withFields(withContext(ctx), flds).Error(e)
}

func Errorf(ctx context.Context, format string, a ...interface{}) {
	withContext(ctx).Errorf(format, normalizeArgs(a)...)
}

// Debug prints debug logs while attempting to JSON dump any non-primitive argument.
func Debug(ctx context.Context, d interface{}, flds ...Fld) {
	withFields(withContext(ctx), flds).Debug(d)
}

// Debugf prints formatted debug logs while attempting to JSON dump any non-primitive argument.
func Debugf(ctx context.Context, format string, a ...interface{}) {
	withContext(ctx).Debugf(format, normalizeArgs(a)...)
}

func Fatal(ctx context.Context, err error) {
	withContext(ctx).Fatal(err)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	withContext(ctx).Fatalf(format, args...)
}

func normalizeArgs(a []interface{}) (n []interface{}) {
	for _, i := range a {
		switch v := i.(type) {
		case string, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8, float32, float64, bool, fmt.Stringer, error:
			n = append(n, v)
		default:
			n = append(n, jsonString(v))
		}
	}
	return
}

func jsonString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func Sync() {
	_ = os.Stderr.Sync()
	_ = os.Stdout.Sync()
}
