package core

import (
	"net/http"

	"github.com/meilihao/wechat/util"
)

// 验证服务器配置, 请使用正式的服务器配置Token
// Method GET
func VerifySetting(token string, w http.ResponseWriter, r *http.Request) bool {
	var isOK bool

	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")
	signature := r.FormValue("signature")
	echostr := r.FormValue("echostr")

	if isOK = util.CheckSignature(signature, token, timestamp, nonce); isOK {
		w.Write([]byte(echostr))
	}

	return isOK
}
