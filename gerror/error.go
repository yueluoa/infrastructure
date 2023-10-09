package gerror

type CodeError interface {
	Code() Code
	Error() string
	Unwrap() error
}
