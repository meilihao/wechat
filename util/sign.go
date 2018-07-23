package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
)

func CheckSignature(signature string, strs ...string) bool {
	return GenerateSign(strs...) == signature
}

// 微信公众号 url 签名
// 微信公众号/企业号 消息体签名.
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421135319
// 注意如果使用的是[公众平台测试账号](https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index),
// 这里的token应该使用测试号的"接口配置信息"里的token, 而且微信服务器的事件推送也是会是推送到"接口配置信息"里的URL上.
func GenerateSign(strs ...string) string {
	sort.Strings(strs)

	tmp := sha1.Sum([]byte(strings.Join(strs, "")))
	return hex.EncodeToString(tmp[:])
}

func DecodeAESKey(raw string) []byte {
	key, err := base64.StdEncoding.DecodeString(raw + "=")
	if err != nil {
		panic(err)
	}

	return key
}

// ciphertext = AES_Encrypt[random(16B) + msg_len(4B) + rawXMLMsg + appId]
func AESDecrypt(key []byte, text string) ([]byte, string, error) {
	data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, "", err
	}
	if len(data) < aes.BlockSize {
		return nil, "", errors.New("ciphertext too short")
	}

	cbc := cipher.NewCBCDecrypter(block, key[:16])
	cbc.CryptBlocks(data, data)

	data = Pkcs7Unpadding(data)

	msgLen := buildMsgLen(data[16:20])

	return data[20 : 20+msgLen], string(data[20+msgLen:]), nil
}
