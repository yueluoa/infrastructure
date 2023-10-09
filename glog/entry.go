package glog

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type Entry struct {
	Ctx     context.Context
	Log     *Log
	Data    []Field // 自定义字段
	Time    time.Time
	Level   Level
	Caller  *runtime.Frame
	Buffer  *bytes.Buffer
	Message string
	err     string
}

func NewEntry(log *Log) *Entry {
	return &Entry{
		Log:  log,
		Data: make([]Field, 0),
	}
}

func (entry *Entry) Bytes() ([]byte, error) {
	return entry.Log.Formatter.Format(entry)
}

// 扩展hook使用
func (entry *Entry) String() (string, error) {
	serialized, err := entry.Bytes()
	str := string(serialized)
	return str, err
}

func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(Field{Key: ErrorKey, Value: err})
}

func (entry *Entry) WithTime(t time.Time) *Entry {
	dataCopy := make([]Field, len(entry.Data))
	for k, v := range entry.Data {
		dataCopy[k] = v
	}
	return &Entry{Ctx: entry.Ctx, Log: entry.Log, Data: dataCopy, Time: t, err: entry.err}
}

func (entry *Entry) WithContext(ctx context.Context) *Entry {
	dataCopy := make([]Field, len(entry.Data))
	for k, v := range entry.Data {
		dataCopy[k] = v
	}
	return &Entry{Ctx: ctx, Log: entry.Log, Data: dataCopy, Time: entry.Time, err: entry.err}
}

func (entry *Entry) WithField(field Field) *Entry {
	return entry.WithFields([]Field{{Key: field.Key, Value: field.Value}})
}

func (entry *Entry) WithFields(fields []Field) *Entry {
	data := make([]Field, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	fieldErr := entry.err
	for k, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch {
			case t.Kind() == reflect.Func, t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Func:
				isErrField = true
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", k)
			if fieldErr != "" {
				fieldErr = entry.err + ", " + tmp
			} else {
				fieldErr = tmp
			}
		} else {
			data[k] = v
		}
	}
	return &Entry{Ctx: entry.Ctx, Log: entry.Log, Data: data, Time: entry.Time, err: fieldErr}
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// 检索第一个非log调用函数的名称
func getCaller() *runtime.Frame {
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// 动态获取包名和最小调用深度
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				logPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogFrames
	})

	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// 如果调用者不是这个包的一部分，就完成了
		if pkg != logPackage {
			return &f
		}
	}

	return nil
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.log(DebugLevel, args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.log(InfoLevel, args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.log(WarnLevel, args...)
}

func (entry *Entry) Warning(args ...interface{}) {
	entry.Warn(args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.log(ErrorLevel, args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.log(FatalLevel, args...)
	entry.Log.Exit()
}

func (entry *Entry) Panic(args ...interface{}) {
	entry.log(PanicLevel, args...)
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.logf(DebugLevel, format, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.logf(InfoLevel, format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.logf(WarnLevel, format, args...)
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.logf(ErrorLevel, format, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.logf(FatalLevel, format, args...)
	entry.Log.Exit()
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.logf(PanicLevel, format, args...)
}

func (entry *Entry) loadLog(level Level, msg string) {
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}

	entry.Level = level
	entry.Message = msg

	entry.Log.mu.Lock()
	reportCaller := entry.Log.ReportCaller
	bufPool := entry.getBufferPool()
	entry.Log.mu.Unlock()

	if reportCaller {
		entry.Caller = getCaller()
	}

	buffer := bufPool.Get()
	defer func() {
		entry.Buffer = nil
		buffer.Reset()
		bufPool.Put(buffer)
	}()
	buffer.Reset()
	entry.Buffer = buffer

	entry.write()

	entry.Buffer = nil

	if level <= PanicLevel {
		panic(entry)
	}
}

func (entry *Entry) write() {
	entry.Log.mu.Lock()
	defer entry.Log.mu.Unlock()
	// read log bytes
	serialized, err := entry.Log.Formatter.Format(entry)
	if err != nil {
		fmt.Printf("read failed, %v\n", err)
		return
	}
	_, err = entry.Log.Out.Write(serialized)
	if err != nil {
		fmt.Printf("failed to write to log, %v\n", err)
	}
}

func (entry *Entry) log(level Level, args ...interface{}) {
	if entry.Log.IsLevelEnabled(level) {
		entry.loadLog(level, fmt.Sprint(args...))
	}
}

func (entry *Entry) logf(level Level, format string, args ...interface{}) {
	if entry.Log.IsLevelEnabled(level) {
		entry.log(level, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) getBufferPool() (pool BufferPool) {
	if entry.Log.BufferPool != nil {
		return entry.Log.BufferPool
	}
	return bufferPool
}
