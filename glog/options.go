package glog

import "io"

type Option interface {
	apply(*Log)
}

type LogOption struct {
	f func(*Log)
}

func (lo *LogOption) apply(log *Log) {
	lo.f(log)
}

func NewLogOption(f func(*Log)) *LogOption {
	return &LogOption{
		f: f,
	}
}

func WithLevel(level Level) Option {
	return NewLogOption(func(log *Log) {
		log.Level = level
	})
}

func WithFormatter(f Formatter) Option {
	return NewLogOption(func(log *Log) {
		log.Formatter = f
	})
}

func WithOutput(output io.Writer) Option {
	return NewLogOption(func(log *Log) {
		log.Out = output
	})
}

func WithReportCaller(reportCaller bool) Option {
	return NewLogOption(func(log *Log) {
		log.ReportCaller = reportCaller
	})
}
