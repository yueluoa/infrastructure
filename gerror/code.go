package gerror

type Code int

const (
	CodeCommon       = 20001
	CodeUnauthorized = 40001
	CodeNotExist     = 40004
)

var (
	CommonError              = &Error{code: CodeCommon, error: new("通用错误")}
	UnauthorizedError        = &Error{code: CodeUnauthorized, error: new("用户未授权")}
	EquipNotExistError       = &Error{code: CodeNotExist, error: new("装备不存在")}
	InscriptionNotExistError = &Error{code: CodeNotExist, error: new("铭纹不存在")}
	GoodsNotExistError       = &Error{code: CodeNotExist, error: new("没有该物品")}
	UserNotExistError        = &Error{code: CodeNotExist, error: new("用户不存在")}
	FriendNotExistError      = &Error{code: CodeNotExist, error: new("好友不存在")}
	AreAlreadyFriendError    = &Error{code: CodeCommon, error: new("你们已经是好友了")}
)
