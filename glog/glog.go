package glog

import (
	"context"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var log *Log
var once sync.Once

type Log struct {
	Out          io.Writer
	Formatter    Formatter
	Level        Level
	entryPool    sync.Pool
	ExitFunc     exitFunc
	BufferPool   BufferPool
	ReportCaller bool // 是否标记调用信息
	mu           sync.Mutex
}

func New(opts ...Option) *Log {
	once.Do(func() {
		log = newLog(opts...)
	})

	return log
}

func newLog(opts ...Option) *Log {
	log := &Log{
		Out:       os.Stderr,
		Formatter: new(TextFormatter),
		Level:     InfoLevel,
		ExitFunc:  os.Exit,
	}
	for _, opt := range opts {
		opt.apply(log)
	}

	return log
}

func (log *Log) newEntry() *Entry {
	entry, ok := log.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(log)
}

func (log *Log) putEntry(entry *Entry) {
	entry.Data = []Field{}
	log.entryPool.Put(entry)
}

func (log *Log) WithField(field Field) *Entry {
	entry := log.newEntry()
	defer log.putEntry(entry)
	return entry.WithField(field)
}

func (log *Log) WithFields(fields []Field) *Entry {
	entry := log.newEntry()
	defer log.putEntry(entry)
	return entry.WithFields(fields)
}

func (log *Log) WithError(err error) *Entry {
	entry := log.newEntry()
	defer log.putEntry(entry)
	return entry.WithError(err)
}

func (log *Log) WithContext(ctx context.Context) *Entry {
	entry := log.newEntry()
	defer log.putEntry(entry)
	return entry.WithContext(ctx)
}

// 修改日志输出时间
func (log *Log) WithTime(t time.Time) *Entry {
	entry := log.newEntry()
	defer log.putEntry(entry)
	return entry.WithTime(t)
}

func (log *Log) log(level Level, args ...interface{}) {
	if log.IsLevelEnabled(level) {
		entry := log.newEntry()
		defer log.putEntry(entry)
		entry.log(level, args...)
	}
}

func (log *Log) Debug(args ...interface{}) {
	log.log(DebugLevel, args...)
}

func (log *Log) Info(args ...interface{}) {
	log.log(InfoLevel, args...)
}

func (log *Log) Warn(args ...interface{}) {
	log.log(WarnLevel, args...)
}

func (log *Log) Warning(args ...interface{}) {
	log.Warn(args...)
}

func (log *Log) Error(args ...interface{}) {
	log.log(ErrorLevel, args...)
}

func (log *Log) Fatal(args ...interface{}) {
	log.log(FatalLevel, args...)
	log.Exit()
}

func (log *Log) Panic(args ...interface{}) {
	log.log(PanicLevel, args...)
}

func (log *Log) logf(level Level, format string, args ...interface{}) {
	if log.IsLevelEnabled(level) {
		entry := log.newEntry()
		defer log.putEntry(entry)
		entry.logf(level, format, args...)
	}
}

func (log *Log) Debugf(format string, args ...interface{}) {
	log.logf(DebugLevel, format, args...)
}

func (log *Log) Infof(format string, args ...interface{}) {
	log.logf(InfoLevel, format, args...)
}

func (log *Log) Warnf(format string, args ...interface{}) {
	log.logf(WarnLevel, format, args...)
}

func (log *Log) Warningf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func (log *Log) Errorf(format string, args ...interface{}) {
	log.logf(ErrorLevel, format, args...)
}

func (log *Log) Fatalf(format string, args ...interface{}) {
	log.logf(FatalLevel, format, args...)
	log.Exit()
}

func (log *Log) Panicf(format string, args ...interface{}) {
	log.logf(PanicLevel, format, args...)
}

func (log *Log) Exit() {
	if log.ExitFunc == nil {
		log.ExitFunc = os.Exit
	}
	log.ExitFunc(exitCode)
}

// 检查日志级别
func (log *Log) IsLevelEnabled(level Level) bool {
	return log.level() >= level
}

func (log *Log) level() Level {
	return Level(atomic.LoadUint32((*uint32)(&log.Level)))
}

func (log *Log) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&log.Level), uint32(level))
}

func (log *Log) GetLevel() Level {
	return log.level()
}

func (log *Log) SetFormatter(formatter Formatter) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.Formatter = formatter
}

func (log *Log) SetOutput(output io.Writer) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.Out = output
}

func (log *Log) SetReportCaller(reportCaller bool) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.ReportCaller = reportCaller
}

func GetLog() *Log {
	return log
}
