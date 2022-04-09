package err_status

import "errors"

type ErrCode int

const (
	Success ErrCode = iota + 2000
	SystemErrCode
)

type ErrMsg struct {
	Code ErrCode `json:"code"`
	Err  error   `json:"err_status"`
}

func New(code ErrCode, err error) *ErrMsg {
	return &ErrMsg{
		Code: code,
		Err:  err,
	}
}
func NewWithMsg(code ErrCode, msg string) *ErrMsg {
	return &ErrMsg{
		Code: code,
		Err:  errors.New(msg),
	}
}
func NewSystemError(err error) *ErrMsg {
	return &ErrMsg{
		Code: SystemErrCode,
		Err:  err,
	}
}

func (e *ErrMsg) Error() string {
	return e.Err.Error()
}
