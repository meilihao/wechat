package util

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"strings"
)

func CheckSignature(signature string, strs ...string) bool {
	return GenerateSign(strs...) == signature
}

// 微信公众号 url 签名
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421135319
// 注意如果使用的是[公众平台测试账号](https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index),
// 这里的token应该使用测试号的"接口配置信息"里的token, 而且微信服务器的事件推送也是会是推送到"接口配置信息"里的URL上.
func GenerateSign(strs ...string) string {
	sort.Strings(strs)

	tmp := sha1.Sum([]byte(strings.Join(strs, "")))
	return hex.EncodeToString(tmp[:])
}
