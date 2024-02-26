package binding

import (
	"git.ejxcgit.com/ejhycommon/go-common/rtjson"
	"git.ejxcgit.com/ejhycommon/go-common/rtjson/msg"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// BusinessValidator 复杂业务逻辑校验器
type BusinessValidator interface {
	// Validate 业务逻辑校验
	Validate(r *http.Request) *rtjson.ServiceReturn
}

// BindThenCheck 请求参数绑定并校验
func BindThenCheck(r *http.Request, req interface{}) *rtjson.ServiceReturn {
	err := httpx.Parse(r, req)
	if err != nil {
		return rtjson.NewErr(msg.RequestParameterErr)
	}

	// 2.复杂业务验证
	if vdt, ok := req.(BusinessValidator); ok {
		return vdt.Validate(r)
	}
	return nil
}
