package glog

const (
	defaultTimestampFormat = "2006-01-02 15:04:05"
	FieldKeyMsg            = "_msg"
	FieldKeyLevel          = "_level"
	FieldKeyTime           = "_time"
	FieldKeyFunc           = "func"
	FieldKeyFile           = "file"
)

type Formatter interface {
	Format(*Entry) ([]byte, error)
}
