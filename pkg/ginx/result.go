package ginx

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success() Result {
	return Result{
		Code: 0,
		Msg:  "请求成功",
	}
}

func InvalidToken() Result {
	return Result{
		Code: 2,
		Msg:  "无效的token",
	}
}

func InvalidParam() Result {
	return Result{
		Code: 4,
		Msg:  "参数错误",
	}
}
func SystemError() Result {
	return Result{
		Code: 5,
		Msg:  "系统错误",
	}
}
