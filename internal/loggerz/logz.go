package loggerz

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	il "github.com/kubex-ecosystem/logz/internal/interfaces"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

var (
	// info      m.Manifest
	debug     bool
	showTrace bool
	logLevel  string
	g         *LoggerImpl
	LoggerG   Logger
	// err       error
)

type Logger interface {
	// Embedding the LogzLogger interface
	il.LogzLogger

	GetLogger() *LoggerImpl
	GetShowTrace() bool
	GetDebug() bool
	SetLogLevel(logLevel string)
	SetShowTrace(showTrace bool)
	SetDebug(d bool)
	Log(logType string, messages ...any)
	ObjLog(obj any, logType string, messages ...any)

	Notice(m ...any)
	Info(m ...any)
	Debug(m ...any)
	Warn(m ...any)
	Error(m ...any)
	Fatal(m ...any)
	Panic(m ...any)
	Success(m ...any)
	Silent(m ...any)
	Answer(m ...any)
}

type LogzImpl struct {
	// Embedding the LogzCore implementation
	*LogzCoreImpl
}

type LoggerImpl struct {
	*LogzImpl

	gLogLevel    LogLevel    // Global log level
	gLogLevelInt LogLevelInt // Global log level
	gShowTrace   bool        // Flag to show trace in logs
	gDebug       bool        // Flag to show debug messages
}

type LogType string
type LogLevel = kbx.LogLevel
type LogLevelInt int

const (
	// LogLevelDebug 0
	LogLevelDebug LogLevelInt = iota
	// LogLevelNotice 1
	LogLevelNotice
	// LogLevelInfo 2
	LogLevelInfo
	// LogLevelSuccess 3
	LogLevelSuccess
	// LogLevelWarn 4
	LogLevelWarn
	// LogLevelError 5
	LogLevelError
	// LogLevelFatal 6
	LogLevelFatal
	// LogLevelPanic 7
	LogLevelPanic
	// LogLevelAnswer 8
	LogLevelAnswer
	// LogLevelSilent 9
	LogLevelSilent
)

func getEnvOrDefault[T string | int | bool](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	} else {
		valInterface := reflect.ValueOf(value)
		if valInterface.Type().ConvertibleTo(reflect.TypeFor[T]()) {
			return valInterface.Convert(reflect.TypeFor[T]()).Interface().(T)
		}
	}
	return defaultValue
}

func (g *LoggerImpl) GetConfig() il.Config {
	if g.LogzImpl == nil || g.VConfig == nil {
		return nil
	}
	// if configImpl, ok := g.VConfig.(il.Config); ok {
	configImpl := g.VConfig
	return configImpl
}
func (g *LoggerImpl) SetConfig(cfg il.Config) {
	if g.LogzImpl == nil {
		return
	}
	g.VConfig = cfg
}
func (g *LoggerImpl) GetLogger() *LoggerImpl              { return g }
func (g *LoggerImpl) GetLogLevel() string                 { return string(g.gLogLevel) }
func (g *LoggerImpl) GetShowTrace() bool                  { return g.gShowTrace }
func (g *LoggerImpl) GetDebug() bool                      { return g.gDebug }
func (g *LoggerImpl) SetLogLevel(logLevel string)         { setLogLevel(logLevel) }
func (g *LoggerImpl) SetShowTrace(showTrace bool)         { g.gShowTrace = showTrace }
func (g *LoggerImpl) SetDebug(d bool)                     { SetDebug(d); g.gDebug = d }
func (g *LoggerImpl) Log(logType string, messages ...any) { Log(logType, messages...) }
func (g *LoggerImpl) ObjLog(obj any, logType string, messages ...any) {
	var T *LoggerImpl
	var ok bool
	if T, ok = obj.(*LoggerImpl); !ok {
		g.ErrorCtx("ObjLog: obj is not of type *Logger", map[string]any{
			"context":  "ObjLog",
			"logType":  logType,
			"object":   obj,
			"msg":      messages,
			"showData": getShowTrace(),
		})
		return
	}
	LogObjLogger[LoggerImpl](T, logType, messages...)
}
func (g *LoggerImpl) Notice(m ...any) {
	if g == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.Log("notice", m...)
}
func (g *LoggerImpl) Info(m ...any) {
	if g == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.Log("info", m...)
}
func (g *LoggerImpl) Debug(m ...any) {
	if g == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.Log("debug", m...)
}
func (g *LoggerImpl) Warn(m ...any) {
	if g == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.Log("warn", m...)
}
func (g *LoggerImpl) Error(m ...any) {
	if g == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.Log("error", m...)
}
func (g *LoggerImpl) Fatal(m ...any) {
	if g == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.Log("fatal", m...)
}
func (g *LoggerImpl) Panic(m ...any) {
	if g == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.Log("fatal", m...)
}
func (g *LoggerImpl) Success(m ...any) {
	if g == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.Log("success", m...)
}
func (g *LoggerImpl) Silent(m ...any) {
	if g == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.Log("silent", m...)
}
func (g *LoggerImpl) Answer(m ...any) {
	if g == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.Log("answer", m...)
}
func (g *LoggerImpl) Init() error {
	if g.Logger == nil {
		_ = GetLoggerZ[*LoggerImpl](nil)
	}
	return nil
}

func NewLoggerzImpl(prefix string) *LoggerImpl {
	return &LoggerImpl{
		LogzImpl:     &LogzImpl{LogzCoreImpl: NewLoggerImpl(prefix)},
		gLogLevel:    LogLevel(kbx.ERROR),
		gLogLevelInt: LogLevelError,
		gShowTrace:   false,
		gDebug:       false,
	}
}
func NewLoggerZ(prefix string) Logger { return NewLoggerzImpl(prefix) }

func GetLogger(prefix *string) Logger {
	if prefix != nil {
		return NewLoggerzImpl(*prefix)
	}
	return g.GetLogger()
}

func GetLoggerZ[T any](obj *T) Logger {
	if g == nil || LoggerG == nil {
		g = &LoggerImpl{
			// Logger:     GetLogger(info.GetBin()),
			gLogLevel:    LogLevel(kbx.INFO),
			gLogLevelInt: LogLevelInfo,
			gShowTrace:   showTrace,
			gDebug:       debug,
		}
		LoggerG = g
	}
	if obj == nil {
		if LoggerG == nil {
			return g
		}
		return LoggerG
	}
	var lgr *LoggerImpl
	if objValueLogger := reflect.ValueOf(obj).Elem().MethodByName("GetLogger"); !objValueLogger.IsValid() {
		// check if is interface, if so, try another approach
		if reflect.TypeOf(obj).Kind() == reflect.Interface {
			if reflect.ValueOf(obj).Elem().Kind() == reflect.Ptr {
				if objValueLogger = reflect.ValueOf(obj).Elem().Elem().MethodByName("GetLogger"); !objValueLogger.IsValid() {
					g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
						"context":  "Log",
						"logType":  "error",
						"object":   obj,
						"msg":      "object does not have a logger field",
						"showData": getShowTrace(),
					})
					return g
				}
			} else {
				g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
					"context":  "Log",
					"logType":  "error",
					"object":   obj,
					"msg":      "object does not have a logger field",
					"showData": getShowTrace(),
				})
				return g
			}
		} else {
			g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
				"context":  "Log",
				"logType":  "error",
				"object":   obj,
				"msg":      "object does not have a logger field",
				"showData": getShowTrace(),
			})
			return g
		}
		lgrC := objValueLogger.Call([]reflect.Value{})
		if len(lgrC) == 0 {
			g.ErrorCtx(fmt.Sprintf("log object (%s) GetLogger method returned no value", reflect.TypeFor[T]()), map[string]any{
				"context":  "Log",
				"logType":  "error",
				"object":   obj,
				"msg":      "object does not have a logger field",
				"showData": getShowTrace(),
			})
			return g
		}
		if lgrC[0].IsNil() {
			lgr = g
		} else {
			if lgrC[0].Type().ConvertibleTo(reflect.TypeFor[*LoggerImpl]()) {
				lgr = lgrC[0].Convert(reflect.TypeFor[*LoggerImpl]()).Interface().(*LoggerImpl)
			} else {
				g.ErrorCtx(fmt.Sprintf("log object (%s) GetLogger method returned invalid type", reflect.TypeFor[T]()), map[string]any{
					"context":  "Log",
					"logType":  "error",
					"object":   obj,
					"msg":      "object does not have a logger field",
					"showData": getShowTrace(),
				})
				return g
			}
		}
	} else {
		lgrC := objValueLogger.Call([]reflect.Value{})
		if len(lgrC) == 0 {
			g.ErrorCtx(fmt.Sprintf("log object (%s) GetLogger method returned no value", reflect.TypeFor[T]()), map[string]any{
				"context":  "Log",
				"logType":  "error",
				"object":   obj,
				"msg":      "object does not have a logger field",
				"showData": getShowTrace(),
			})
			return g
		}
		if lgrC[0].IsNil() {
			lgr = g
		} else {
			if lgrC[0].Type().ConvertibleTo(reflect.TypeFor[*LoggerImpl]()) {
				lgr = lgrC[0].Convert(reflect.TypeFor[*LoggerImpl]()).Interface().(*LoggerImpl)
			} else {
				g.ErrorCtx(fmt.Sprintf("log object (%s) GetLogger method returned invalid type", reflect.TypeFor[T]()), map[string]any{
					"context":  "Log",
					"logType":  "error",
					"object":   obj,
					"msg":      "object does not have a logger field",
					"showData": getShowTrace(),
				})
				return g
			}
		}
	}
	if lgr == nil {
		g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
			"context":  "Log",
			"logType":  "error",
			"object":   obj,
			"msg":      "object does not have a logger field",
			"showData": getShowTrace(),
		})
		return LoggerG
	}
	return &LoggerImpl{
		LogzImpl:   lgr.LogzImpl,
		gLogLevel:  g.gLogLevel,
		gShowTrace: g.gShowTrace,
		gDebug:     g.gDebug,
	}
}

func Log(logType string, messages ...any) {
	funcName, line, file := getFuncNameMessage(g)
	fullMessage := getFullMessage(messages...)
	logType = strings.ToUpper(logType)
	ctxMessageMap := getCtxMessageMap(logType, funcName, file, line)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			ctxMessageMap["logType"] = logType
			logging(g, lType, fullMessage, ctxMessageMap)
		} else {
			g.ErrorCtx(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		logging(g, LogType(logType), fullMessage, ctxMessageMap)
	}
}

func logging(lgr *LoggerImpl, lType LogType, fullMessage string, ctxMessageMap map[string]any) {
	lt := strings.ToLower(string(lType))
	if _, exist := ctxMessageMap["showData"]; !exist {
		ctxMessageMap["showData"] = getShowTrace()
	}
	if willPrintLog(lt) {
		switch lType {
		case LogType(kbx.INFO):
			lgr.InfoCtx(fullMessage, ctxMessageMap)
		case LogType(kbx.DEBUG):
			lgr.DebugCtx(fullMessage, ctxMessageMap)
		case LogType(kbx.ERROR):
			lgr.ErrorCtx(fullMessage, ctxMessageMap)
		case LogType(kbx.WARN):
			lgr.WarnCtx(fullMessage, ctxMessageMap)
		case LogType(kbx.NOTICE):
			lgr.NoticeCtx(fullMessage, ctxMessageMap)
		case LogType(kbx.SUCCESS):
			lgr.SuccessCtx(fullMessage, ctxMessageMap)
		case LogType(kbx.FATAL):
			lgr.FatalCtx(fullMessage, ctxMessageMap)
		case LogType(kbx.PANIC):
			lgr.FatalCtx(fullMessage, ctxMessageMap)
		case LogType(kbx.SILENT):
			lgr.SilentCtx(fullMessage, ctxMessageMap)
		case LogType(kbx.ANSWER):
			lgr.AnswerCtx(fullMessage, ctxMessageMap)
		default:
			lgr.InfoCtx(fullMessage, ctxMessageMap)
		}
	} else {
		ctxMessageMap["msg"] = fullMessage
		ctxMessageMap["showData"] = false
		lgr.DebugCtx(ctxMessageMap["msg"].(string), ctxMessageMap)
	}
}

func SetDebugMode(debug bool) {
	SetDebug(debug)
}

func SetLogLevel(level string) {
	LoggerG.SetLogLevel(level)
}

func SetLogTrace(enable bool) {
	LoggerG.SetShowTrace(enable)
}

func SetLogger(logger *LoggerImpl) {
	// gl.SetLogger(logger)
	// TODO: Implement this function properly
}

func setLogLevel(logLevel string) {
	if g == nil || LoggerG == nil {
		// _ = GetLoggerZ[LoggerImpl](nil)
		g = NewLoggerZ("").(*LoggerImpl)
		LoggerG = g
	}
	switch strings.ToLower(logLevel) {
	case "debug":
		g.gLogLevel = kbx.DEBUG
		g.gLogLevelInt = LogLevelDebug
		g.SetLogLevel("debug")
	case "info":
		g.gLogLevel = kbx.INFO
		g.gLogLevelInt = LogLevelInfo
		g.SetLevel("info")
	case "warn":
		g.gLogLevel = kbx.WARN
		g.gLogLevelInt = LogLevelWarn
		g.SetLevel("warn")
	case "error":
		g.gLogLevel = kbx.ERROR
		g.gLogLevelInt = LogLevelError
		g.SetLevel("error")
	case "fatal":
		g.gLogLevel = kbx.FATAL
		g.gLogLevelInt = LogLevelFatal
		g.SetLevel("fatal")
	case "panic":
		g.gLogLevel = kbx.PANIC
		g.gLogLevelInt = LogLevelPanic
		g.SetLevel("panic")
	case "notice":
		g.gLogLevel = kbx.NOTICE
		g.gLogLevelInt = LogLevelNotice
		g.SetLevel("notice")
	case "success":
		g.gLogLevel = kbx.SUCCESS
		g.gLogLevelInt = LogLevelSuccess
		g.SetLevel("success")
	case "silent":
		g.gLogLevel = kbx.SILENT
		g.gLogLevelInt = LogLevelSilent
		g.SetLevel("silent")
	case "answer":
		g.gLogLevel = kbx.ANSWER
		g.gLogLevelInt = LogLevelAnswer
		g.SetLevel("answer")
	default:
		logLevel = "info"
		g.gLogLevel = kbx.INFO
		g.gLogLevelInt = LogLevelInfo
		g.SetLevel(logLevel)
	}
}
func getShowTrace() bool {
	if debug {
		showTrace = true
		return true
	} else {
		if !showTrace {
			return false
		} else {
			return true
		}
	}
}
func willPrintLog(logType string) bool {
	if debug {
		return true
	} else {
		lTypeInt := LogLevelError
		switch strings.ToLower(logType) {
		case "debug":
			lTypeInt = LogLevelDebug
		case "info":
			lTypeInt = LogLevelInfo
		case "warn":
			lTypeInt = LogLevelWarn
		case "error":
			lTypeInt = LogLevelError
		case "notice":
			lTypeInt = LogLevelNotice
		case "success":
			lTypeInt = LogLevelSuccess
		case "fatal":
			lTypeInt = LogLevelFatal
		case "panic":
			lTypeInt = LogLevelPanic
		case "silent":
			lTypeInt = LogLevelSilent
		case "answer":
			lTypeInt = LogLevelAnswer
		default:
			lTypeInt = LogLevelError
		}

		return lTypeInt >= g.gLogLevelInt
	}
}
func getCtxMessageMap(logType, funcName, file string, line int) map[string]any {
	ctxMessageMap := map[string]any{
		"context":   funcName,
		"file":      file,
		"line":      line,
		"logType":   logType,
		"timestamp": time.Now().Format(time.RFC3339),
		// "version":   info.GetVersion(),
	}
	if !debug && !showTrace {
		ctxMessageMap["showData"] = false
	} else {
		ctxMessageMap["showData"] = getShowTrace()
	}
	// if info != nil {
	// 	ctxMessageMap["appName"] = info.GetName()
	// 	ctxMessageMap["bin"] = info.GetBin()
	// 	ctxMessageMap["version"] = info.GetVersion()
	// }
	return ctxMessageMap
}
func getFuncNameMessage(lgr Logger) (string, int, string) {
	if getShowTrace() {
		pc, file, line, ok := runtime.Caller(4)
		if !ok {
			lgr.ErrorCtx("Log: unable to get caller information", nil)
			return "", 0, ""
		}
		funcName := runtime.FuncForPC(pc).Name()
		if strings.Contains(funcName, "LogObjLogger") {
			pc, file, line, ok = runtime.Caller(4)
			if !ok {
				lgr.ErrorCtx("Log: unable to get caller information", nil)
				return "", 0, ""
			}
			funcName = runtime.FuncForPC(pc).Name()
		}
		return funcName, line, file
	}
	return "", 0, ""
}
func getFullMessage(messages ...any) string {
	fullMessage := ""
	for _, msg := range messages {
		if msg != nil {
			if str, ok := msg.(string); ok {
				fullMessage += str + " "
			} else {
				fullMessage += fmt.Sprintf("%v ", msg)
			}
		}
	}
	return strings.TrimSuffix(
		strings.TrimPrefix(
			strings.TrimSpace(fullMessage),
			" ",
		),
		" ",
	)
}

func SetDebug(d bool) {
	if g == nil || LoggerG == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	g.gDebug = d
	if d {
		showTrace = true
		debug = true
		g.SetLevel("debug")
	} else {
		switch g.gLogLevelInt {
		case LogLevelDebug:
			g.SetLevel("debug")
		case LogLevelInfo:
			g.SetLevel("info")
		case LogLevelWarn:
			g.SetLevel("warn")
		case LogLevelError:
			g.SetLevel("error")
		case LogLevelFatal:
			g.SetLevel("fatal")
		case LogLevelPanic:
			g.SetLevel("panic")
		case LogLevelNotice:
			g.SetLevel("notice")
		case LogLevelSuccess:
			g.SetLevel("success")
		case LogLevelSilent:
			g.SetLevel("silent")
		case LogLevelAnswer:
			g.SetLevel("answer")
		default:
			g.SetLevel("info")
		}
	}
}
func LogObjLogger[T any](obj *T, logType string, messages ...any) {
	defer func() {
		if r := recover(); r != nil {
			if g == nil || LoggerG == nil {
				_ = GetLoggerZ[LoggerImpl](nil)
			}
			g.ErrorCtx(fmt.Sprintf("LogObjLogger panic: %v", r), map[string]any{
				"context":  "LogObjLogger",
				"logType":  logType,
				"object":   obj,
				"msg":      messages,
				"showData": getShowTrace(),
			})
		}
	}()
	if g == nil || LoggerG == nil {
		_ = GetLoggerZ[LoggerImpl](nil)
	}
	if obj == nil {
		g.ErrorCtx("LogObjLogger: obj is nil", map[string]any{
			"context":  "LogObjLogger",
			"logType":  logType,
			"object":   obj,
			"msg":      messages,
			"showData": getShowTrace(),
		})
		return
	}

	lgr := GetLoggerZ(obj)
	if lgr == nil {
		g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
			"context":  "Log",
			"logType":  logType,
			"object":   obj,
			"msg":      messages,
			"showData": getShowTrace(),
		})
		return
	}

	fullMessage := getFullMessage(messages...)
	logType = strings.ToLower(logType)
	funcName, line, file := getFuncNameMessage(lgr.GetLogger())

	ctxMessageMap := getCtxMessageMap(logType, funcName, file, line)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			logging(lgr.GetLogger(), lType, fullMessage, ctxMessageMap)
		} else {
			lgr.GetLogger().ErrorCtx(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		lgr.GetLogger().InfoCtx(fullMessage, ctxMessageMap)
	}
}
