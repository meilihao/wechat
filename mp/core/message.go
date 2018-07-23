package core

import (
	"encoding/xml"

	"github.com/meilihao/wechat/util"
)

type MsgCryptor struct {
	AppId  string
	Token  string
	AESKey []byte
}

func NewMsgCryptor(appId, token, encodingAESKey string) *MsgCryptor {
	return &MsgCryptor{
		Token:  token,
		AESKey: util.DecodeAESKey(encodingAESKey),
		AppId:  appId,
	}
}

// 来自微信服务器的加密消息
type CipherRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	EncryptedMsg string   `xml:"Encrypt"` // base64 std encoded
}

// 发送给微信服务器的消息
type CipherResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	EncryptedMsg string   `xml:"Encrypt"` // base64 std encoded
	MsgSignature string   `xml:"MsgSignature"`
	TimeStamp    string   `xml:"TimeStamp"`
	Nonce        string   `xml:"Nonce"`
}

// MsgType 基本消息类型
type MsgType string

const (
	// MsgTypeEvent 事件推送消息
	MsgTypeEvent = "event"
)

// Event 事件类型
type Event string

const (
	// EventSubscribe 订阅
	EventSubscribe Event = "subscribe"
	// EventUnsubscribe 取消订阅
	EventUnsubscribe = "unsubscribe"
)

type MixedMsg struct {
	XMLName xml.Name `xml:"xml"`
	MsgHeader
	Event Event `xml:"Event" json:"Event"`
}

// 微信服务器推送过来的消息(事件)的通用消息头.
type MsgHeader struct {
	ToUserName   string  `xml:"ToUserName"   json:"ToUserName"`
	FromUserName string  `xml:"FromUserName" json:"FromUserName"`
	CreateTime   int64   `xml:"CreateTime"   json:"CreateTime"`
	MsgType      MsgType `xml:"MsgType"      json:"MsgType"`
}
