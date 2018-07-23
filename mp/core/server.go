package core

import (
	"bytes"
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
func (or *MsgCryptor) Decrypt(raw []byte, us url.Values) ([]byte, error) {
	var tmpBody CipherRequestBody

	err := xml.Unmarshal(raw, &tmpBody)
	if err != nil {
		return nil, err
	}

	if !util.CheckSignature(us.Get("msg_signature"), or.Token, us.Get("timestamp"), us.Get("nonce"), tmpBody.EncryptedMsg) {
		return nil, errors.New("invalid msg_signature")
	}

	data, tmpAppId, err := util.AESDecrypt(or.AESKey, tmpBody.EncryptedMsg)
	if err != nil {
		return nil, err
	}

	if tmpAppId != or.AppId {
		return nil, errors.New("get other appid from msg")
	}

	return data, nil
}

// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1419318482&token=&lang=zh_CN
func (or *MsgCryptor) Encrypt(raw []byte, timestamp, nonce string) *bytes.Buffer {
	encryptedMsg := util.AESEecrypt(or.AESKey, raw, or.AppId)
	msgSignature := util.GenerateSign(or.Token, timestamp, nonce, string(encryptedMsg))

	buf := bytes.NewBuffer(make([]byte, 0, len(encryptedMsg)))

	buf.WriteString("<xml><Encrypt><![CDATA[")
	buf.Write(encryptedMsg)
	buf.WriteString("]]></Encrypt>")
	buf.WriteString("<MsgSignature><![CDATA[")
	buf.WriteString(msgSignature)
	buf.WriteString("]]></MsgSignature>")
	buf.WriteString("<TimeStamp><![CDATA[")
	buf.WriteString(timestamp)
	buf.WriteString("]]></TimeStamp>")
	buf.WriteString("<Nonce><![CDATA[")
	buf.WriteString(nonce)
	buf.WriteString("]]></Nonce></xml>")

	return buf
}
