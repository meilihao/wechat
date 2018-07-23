package core

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"

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

// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1434696670
// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1419318479&token=&lang=zh_CN
func DecryptMsg(raw, aesKey []byte, appId, token string, us url.Values) ([]byte, error) {
	var tmpBody CipherBody

	err := xml.Unmarshal(raw, &tmpBody)
	if err != nil {
		return nil, err
	}

	if !util.CheckSignature(us.Get("msg_signature"), token, us.Get("timestamp"), us.Get("nonce"), tmpBody.EncryptedMsg) {
		return nil, errors.New("invalid msg_signature")
	}

	data, tmpAppId, err := util.AESDecrypt(aesKey, tmpBody.EncryptedMsg)
	if err != nil {
		return nil, err
	}

	if tmpAppId != appId {
		return nil, errors.New("get other appid from msg")
	}

	return data, nil
}

type CipherBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	EncryptedMsg string   `xml:"Encrypt"` // base64 std encoded
}

type MixedMsg struct {
	XMLName xml.Name `xml:"xml"`
	MsgHeader
	Event string `xml:"Event" json:"Event"`
}

// 微信服务器推送过来的消息(事件)的通用消息头.
type MsgHeader struct {
	ToUserName   string `xml:"ToUserName"   json:"ToUserName"`
	FromUserName string `xml:"FromUserName" json:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"   json:"CreateTime"`
	MsgType      string `xml:"MsgType"      json:"MsgType"`
}
