package global

import (
	"fmt"
	"net/http"
)

type CheckedParams struct {
	Func    int
	Checked bool
	Self    bool
	Pass    []byte
	V       map[string]interface{}
}

type Handler interface {
	Init()
	Handle(checkedParams *CheckedParams) (response interface{}, err *GotrixError)
	GetPass(token []byte) (response interface{}, err *GotrixError)
	GetSession(token []byte) (response interface{}, err *GotrixError)
}

type Checker interface {
	Check(r *http.Request, handler Handler) (checkedParams *CheckedParams, err *GotrixError)
}

type Asyncer interface {
	Init()
	Async(checkedParams *CheckedParams)
}

type GotrixError struct {
	Status int
	Msg    string
}

func NewGotrixError(gotrixError *GotrixError, a ...interface{}) *GotrixError {
	return &GotrixError{Status: gotrixError.Status, Msg: fmt.Sprintf(gotrixError.Msg, a...)}
}

var BLANK_ERROR *GotrixError = &GotrixError{Status: 1000, Msg: "%s"}
var NO_BODY_ERROR *GotrixError = &GotrixError{Status: 1001, Msg: "HTTP请求时Body为空"}
var READ_BODY_ERROR *GotrixError = &GotrixError{Status: 1002, Msg: "读取HTTP请求的Body时出错"}
var BODY_SCHEME_ERROR *GotrixError = &GotrixError{Status: 1003, Msg: "请求格式有误"}
var PASSWORD_ERROR *GotrixError = &GotrixError{Status: 1004, Msg: "用户密码错误"}
var FUNC_NOT_EXISTS *GotrixError = &GotrixError{Status: 1005, Msg: "接口[%d]不存在"}
var READ_PARAM_ERROR *GotrixError = &GotrixError{Status: 1006, Msg: "读取请求参数时出现异常，参数为空"}
var FUNC_PARAM_MUST *GotrixError = &GotrixError{Status: 1007, Msg: "必要参数[func]没有设值"}
var FUNC_PARAM_ERROR *GotrixError = &GotrixError{Status: 1008, Msg: "必要参数[func]是整形参数，请传入正确的参数值"}
var RETURN_DATE_ECNRYPT_ERROR *GotrixError = &GotrixError{Status: 1009, Msg: "返回数据在加密时出现异常"}
var USER_NOT_EXISTES *GotrixError = &GotrixError{Status: 1010, Msg: "该用户不存在"}
var USER_SESSION_NOT_EXISTES *GotrixError = &GotrixError{Status: 1011, Msg: "会话已过期，请重新登陆"}
var FUNC_SELF_ERROR *GotrixError = &GotrixError{Status: 1012, Msg: "功能号的自解密状态与配置不一致，不允许继续执行"}
var NOT_SUPPORT_CONTENT_TYPE *GotrixError = &GotrixError{Status: 1013, Msg: "不支持的Content-type"}
var FUNC_PRIVATE_ERROR *GotrixError = &GotrixError{Status: 1014, Msg: "该方法是私有方法，外部不能调用"}

var PARAM_MUST_EXISTS *GotrixError = &GotrixError{Status: 2001, Msg: "必要参数[%s]没有设值"}
var PARAM_LENGTH_ERROR *GotrixError = &GotrixError{Status: 2002, Msg: "参数[%s]的长度必须在[%d]和[%d]之间"}
var PARAM_NOT_INTEGER *GotrixError = &GotrixError{Status: 2003, Msg: "整形参数[%s]在转换成整数时出错"}
var PARAM_VALUE_ERROR *GotrixError = &GotrixError{Status: 2004, Msg: "整形参数[%s]的大小必须在[%d]和[%d]之间"}
var PARAM_NOT_BOOLEAN *GotrixError = &GotrixError{Status: 2005, Msg: "布尔参数[%s]在转换成布尔时出错"}
var REDIS_CONNECT_ERROR *GotrixError = &GotrixError{Status: 2006, Msg: "后台数据库连接失败"}
var REDIS_EXEC_ERROR *GotrixError = &GotrixError{Status: 2007, Msg: "后台数据库在执行时出现异常"}
var STRING_PARAM_ERROR *GotrixError = &GotrixError{Status: 2008, Msg: "参数[%s]是字符串型参数，请传入正确的参数值"}
var INT_PARAM_ERROR *GotrixError = &GotrixError{Status: 2009, Msg: "参数[%s]是整型参数，请传入正确的参数值"}
var BOOL_PARAM_ERROR *GotrixError = &GotrixError{Status: 2010, Msg: "参数[%s]是布尔型参数，请传入正确的参数值"}
var JOB_FUNC_NOT_FOUND *GotrixError = &GotrixError{Status: 2011, Msg: "单一功能函数[%s]不存在，请联系开发人员检查配置文件"}
var REDIS_HANDLE_JSON_ERROR *GotrixError = &GotrixError{Status: 2012, Msg: "redisHandle转换成JSON时出现异常"}
var FUNC_JGET_NIL_ERROR *GotrixError = &GotrixError{Status: 2013, Msg: "执行Json操作时，JSON为空"}
var FUNC_JSET_NIL_ERROR *GotrixError = &GotrixError{Status: 2014, Msg: "执行Jset操作时，JSON为空"}
var SQLHANDLE_CONNECT_ERROR *GotrixError = &GotrixError{Status: 2015, Msg: "执行SqlHandle，连接数据库出错"}
var SQLHANDLE_PREPARE_ERROR *GotrixError = &GotrixError{Status: 2016, Msg: "执行SqlHandle，准备SQL语句时出错"}
var SQLHANDLE_QUERY_ERROR *GotrixError = &GotrixError{Status: 2017, Msg: "执行SqlHandle，查询时出错"}
var SQLHANDLE_EXEC_ERROR *GotrixError = &GotrixError{Status: 2018, Msg: "执行SqlHandle，运行时出错"}
var SQLHANDLE_COLUMNS_ERROR *GotrixError = &GotrixError{Status: 2019, Msg: "执行SqlHandle，读取列名称时出错"}
var SQLHANDLE_SCAN_ERROR *GotrixError = &GotrixError{Status: 2020, Msg: "执行SqlHandle,读取行数据时出错"}
var SQLHANDLE_CLOSE_ERROR *GotrixError = &GotrixError{Status: 2021, Msg: "执行SqlHandle，关闭行数据时出错"}
var ARRAY_PARAM_ERROR *GotrixError = &GotrixError{Status: 2022, Msg: "参数[%s]是数组类型参数，请传入正确的参数值"}
var HTTPHANDLE_HTTP_GET_ERROR *GotrixError = &GotrixError{Status: 2050, Msg: "执行httpHandle时，向url请求时出错"}
var HTTPHANDLE_HTTP_READ_BODY *GotrixError = &GotrixError{Status: 2051, Msg: "执行httpHandle时，读取url返回内容时出错"}
var HTTPHANDLE_ANALYZE_ERROR *GotrixError = &GotrixError{Status: 2052, Msg: "执行httpHandle时，解析返回内容时出错"}
var JSON_TO_STRING_ERROR *GotrixError = &GotrixError{Status: 2060, Msg: "执行JsonToString时，解析为json字符串时出错"}

var INTERNAL_ERROR *GotrixError = &GotrixError{Status: 9999, Msg: "内部错误"}