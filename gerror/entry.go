package gerror

type entry struct {
	cause error
	msg   string
}

func newEntry(err error, msg string) error {
	return &entry{
		cause: err,
		msg:   msg,
	}
}

func (entry *entry) Error() string { return entry.msg + ": " + entry.cause.Error() }
func (entry *entry) Unwrap() error { return entry.cause }
func (entry *entry) Cause() error  { return entry.cause }

func new(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
