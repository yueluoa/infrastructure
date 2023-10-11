package gerror

type Code int

const (
	CodeCommon       = 20001
	CodeUnauthorized = 40001
	CodeNotExist     = 40004
)

var (
	CommonError       = &Error{code: CodeCommon, error: new("通用错误")}
	UnauthorizedError = &Error{code: CodeUnauthorized, error: new("用户未授权")}
	DataNotExistError = &Error{code: CodeNotExist, error: new("数据不存在")}
)
