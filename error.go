package n_error

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// PROJECT_ERROR 此專案訂定的，代表未設定正確的 http status
const (
	StatusProjectError   = 512
	StatusDBError        = 513
	StatusExServiceError = 514
)

type Tracer interface {
	error
	Contact() Tracer
}

var _ Tracer = (*ProjectError)(nil)

type ProjectError struct {
	HttpStatus int    `json:"http_status"`
	Code       uint   `json:"code"`
	Msg        string `json:"msg"`
	Params     string `json:"params"`
	Tracer     string `json:"tracer"`
	Err        error  `json:"err"`
	ServiceID  uint   `json:"service_id"`
}

// New 回傳 cash-cow 專屬的錯誤類型 ProjectError
//
// httpStatus 限定 使用 net/http 裡的狀態
//
// code 為此專案唯一值，表示案發地點
//
// msg 錯誤訊息
//
// params 發生錯誤函式，所帶入的參數，可選
func New(httpStatus int, code uint, msg string, err error, params ...any) *ProjectError {
	if httpStatus < 100 || httpStatus >= 600 {
		httpStatus = StatusProjectError // 此專案訂定的，代表未設定正確的 http status
	}

	p := make([]string, len(params))
	for i, param := range params {
		p[i] = fmt.Sprintf("%#v", param)
	}

	file, line := caller()

	return &ProjectError{
		Code:       code,
		Msg:        msg,
		Params:     strings.Join(p, ":"),
		HttpStatus: httpStatus,
		Tracer:     fmt.Sprintf("(%s:%d)", file, line),
		Err:        err,
	}
}

func NewDBErr(code uint, err error, params ...any) *ProjectError {
	return New(StatusDBError, code, "", err, params)
}

func NewInternalErr(code uint, err error, params ...any) *ProjectError {
	return New(http.StatusInternalServerError, code, "", err, params)
}

func NewExService(code uint, serviceID uint, err error, params ...any) *ProjectError {
	projectErr := New(StatusExServiceError, code, "", err, params)
	projectErr.ServiceID = serviceID
	return projectErr
}

func (e *ProjectError) Error() string {
	if e.HttpStatus == StatusProjectError {
		return "http status 未正確設定"
	}

	return fmt.Sprintf("(%d) %s", e.Code, e.Msg)
}

// Contact 接觸過發生錯誤的函式
func (e *ProjectError) Contact() Tracer {
	file, line := caller()
	e.Tracer = fmt.Sprintf("(%s:%d)->%s", file, line, e.Tracer)
	return e
}

// caller 取得呼叫此函數的檔案位置
func caller() (string, int) {
	_, file, line, _ := runtime.Caller(2)

	if paths := strings.Split(file, "/"); len(paths) > 2 {
		return strings.Join(paths[len(paths)-2:], "/"), line
	}

	short := file

	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	file = short

	return file, line
}
