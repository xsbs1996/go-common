package rtjson

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strings"

	"git.ejxcgit.com/ejhycommon/go-common/rtjson/msg"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type RespMessageApp struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Result  interface{} `json:"result"`
}

type ServiceReturn struct {
	Err      error
	Code     int
	Info     interface{}
	Variable map[string]string
	isRpc    bool
}

// NewOK 成功
func NewOK(info interface{}) *ServiceReturn {
	return newServiceReturn(nil, 0, info)
}

// NewErr 普通失败
func NewErr(code int) *ServiceReturn {
	return newServiceReturn(nil, code, nil)
}

// NewHttpErr http 400/500失败返回
func NewHttpErr(err error) *ServiceReturn {
	logx.Errorf("http 500 error:%v", err)
	return newServiceReturn(errors.New("internal error"), msg.HttpErr, nil)
}

// NewRpcErr Rpc失败返回
func NewRpcErr(code int, message string) *ServiceReturn {
	return &ServiceReturn{
		Err:   errors.New(message),
		Code:  code,
		isRpc: true,
	}
}

// NewVarErr 失败,带有变量信息返回
func NewVarErr(code int, variable map[string]string) *ServiceReturn {
	var rt ServiceReturn
	rt.Info = nil
	rt.Code = code
	rt.Variable = variable
	return &rt
}

func newServiceReturn(err error, code int, info interface{}) *ServiceReturn {
	var rt ServiceReturn
	if err != nil {
		rt.Err = err
	}
	if code > 0 {
		rt.Code = code
	}
	if info != nil {
		rt.Info = info
	}
	rt.Variable = nil
	return &rt
}

// JSON 返回
func JSON(w http.ResponseWriter, r *http.Request, rt *ServiceReturn) {
	// rpc错误直接返回
	if rt.isRpc {
		resp := RespMessageApp{
			Message: rt.Err.Error(),
			Code:    rt.Code,
			Result:  nil,
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
		return
	}

	// http400/500以上错误，为了保证熔断降级等业务正常运行
	if rt.Err != nil && rt.Code == msg.HttpErr {
		httpx.ErrorCtx(r.Context(), w, rt.Err)
		return
	}
	if rt.Code > 0 {
		reqAppError(w, r, rt.Code, rt.Variable)
		return
	}
	reqAppSuccess(w, r, rt.Info)
	return
}

// ReqAppSuccess 成功
func reqAppSuccess(w http.ResponseWriter, r *http.Request, result interface{}) {
	resp := RespMessageApp{
		Message: "OK",
		Code:    0,
		Result:  result,
	}
	httpx.OkJsonCtx(r.Context(), w, resp)
}

// ReqAppError 失败
func reqAppError(w http.ResponseWriter, r *http.Request, code int, variable map[string]string) {
	lang := r.Header.Get("lang")
	resp := RespMessageApp{
		Message: msg.GetLangError(lang, code),
		Code:    code,
	}

	if variable != nil {
		for k, v := range variable {
			resp.Message = strings.Replace(resp.Message, fmt.Sprintf("{{$%s}}", k), v, 1)
		}
	}

	httpx.OkJsonCtx(r.Context(), w, resp)
}
