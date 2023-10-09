package glog

import "sync"

const (
	maximumCallerDepth int = 25
	knownLogFrames     int = 4
)

var exitCode = 1

var ErrorKey = "error"

type exitFunc func(int)

var (
	//
	logPackage string
	// 跟踪报告调用方法时在调用堆栈中的位置
	minimumCallerDepth int
	callerInitOnce     sync.Once
)

type Field struct {
	Key   string
	Value interface{}
}

func init() {
	minimumCallerDepth = 1
}
