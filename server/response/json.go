package response

import "encoding/json"

type Response struct {
	Msg  string      `response:"msg"`
	Code int         `response:"code"`
	Info interface{} `response:"info"`
}

func Success(info interface{}) []byte {
	result, _ := json.Marshal(Response{Msg: "ok", Code: 0, Info: info})
	return result

}
func Fail(code int, msg string) []byte {
	result, _ := json.Marshal(Response{Msg: msg, Code: code, Info: nil})
	return result
}
