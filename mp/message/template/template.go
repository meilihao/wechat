// define from https://gopkg.in/chanxuehong/wechat.v2
package template

import (
	"net/url"

	"github.com/meilihao/wechat/util"
)

type TemplateMessage struct {
	ToUser      string       `json:"touser"`                // 必须, 接受者OpenID
	TemplateId  string       `json:"template_id"`           // 必须, 模版ID
	URL         string       `json:"url,omitempty"`         // 可选, 用户点击后跳转的URL, 该URL必须处于开发者在公众平台网站中设置的域中
	MiniProgram *MiniProgram `json:"miniprogram,omitempty"` // 可选, 跳小程序所需数据，不需跳小程序可不用传该数据
	Data        interface{}  `json:"data"`                  // 必须, 模板数据, encoding/json.Marshal 后的数据.
}

type MiniProgram struct {
	AppId    string `json:"appid"`    // 必选; 所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系）
	PagePath string `json:"pagepath"` // 必选; 所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系）
}

// 模版内某个 .DATA 的值
type DataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"` // 模板内容字体颜色，不填默认为黑色
}

// 发送模板消息
// 必须先关注公众号
func Send(accessToken string, msg *TemplateMessage) (msgid int64, err error) {
	baseURL := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" +
		url.QueryEscape(accessToken)

	var result struct {
		util.Error
		MsgId int64 `json:"msgid"`
	}

	if err = util.PostJSON(baseURL, msg, &result); err != nil {
		return
	}
	if result.ErrCode != util.ErrCodeOK {
		err = result.Error

		return
	}

	msgid = result.MsgId

	return
}
