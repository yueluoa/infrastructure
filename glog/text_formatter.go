package glog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type TextFormatter struct {
	pool         sync.Pool
	LevelText    string
	TimeText     string
	DisableColor bool // 是否禁用颜色，默认不禁用
	CallerFrame  func(*runtime.Frame) (function string, file string)
}

func (tf *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var funcVal, fileVal string
	if entry.Caller != nil {
		if tf.CallerFrame != nil {
			funcVal, fileVal = tf.CallerFrame(entry.Caller)
		} else {
			funcVal = entry.Caller.Function
			fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		}
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	tf.TimeText = entry.Time.Format(defaultTimestampFormat)
	tf.LevelText = " [" + strings.ToUpper(entry.Level.String()) + "]"

	if tf.isColored(entry.Log.Out) {
		tf.printColored(b, entry, funcVal, fileVal)
	} else {
		tf.appendValue(b, tf.TimeText)
		tf.appendValue(b, tf.LevelText)
		if funcVal != "" {
			tf.appendKeyValue(b, FieldKeyFunc, funcVal)
		}
		if fileVal != "" {
			tf.appendKeyValue(b, FieldKeyFile, fileVal)
		}
		for _, v := range entry.Data {
			tf.appendKeyValue(b, v.Key, v.Value)
		}
		tf.appendKeyValue(b, FieldKeyMsg, entry.Message)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (tf *TextFormatter) isColored(w io.Writer) bool {
	var isColored bool

	if w == os.Stderr {
		isColored = true
	}

	return isColored && !tf.DisableColor
}

func (tf *TextFormatter) printColored(b *bytes.Buffer, entry *Entry, funcVal, fileVal string) {
	var levelColor int
	switch entry.Level {
	case DebugLevel:
		levelColor = gray
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = red
	case InfoLevel:
		levelColor = blue
	default:
		levelColor = blue
	}

	tf.appendValue(b, tf.TimeText)
	baseText := fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, tf.LevelText)
	tf.appendValue(b, baseText)
	if funcVal != "" {
		tf.appendKeyValue(b, FieldKeyFunc, funcVal)
	}
	if fileVal != "" {
		tf.appendKeyValue(b, FieldKeyFile, fileVal)
	}
	for _, v := range entry.Data {
		tf.appendKeyValue(b, v.Key, v.Value)
	}
	tf.appendKeyValue(b, FieldKeyMsg, entry.Message)
}

func (tf *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteString("= ")
	tf.appendValue(b, value)
}

func (tf *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	b.WriteString(stringVal)
}
