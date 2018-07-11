package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/meilihao/wechat/util"
)

type Client struct {
	AppID, AppSecret string
	Endpoint

	Token      *Token
	HttpClient *http.Client
}

// 网页授权
func NewClient(appID, appSecret string, endpoint Endpoint, httpClient *http.Client) *Client {
	c := &Client{
		AppID:      appID,
		AppSecret:  appSecret,
		Endpoint:   endpoint,
		HttpClient: httpClient,
	}

	if c.HttpClient == nil {
		c.HttpClient = http.DefaultClient
	}

	return c
}

// 用于获取用户基本信息, 无调用次数上限, 但请立即使用
// 此access_token与openid绑定
// 此access_token与基础支持的access_token不同
func (c *Client) Exchange(code string) (*Token, error) {
	resp, err := c.HttpClient.Get(c.exchangeURL(code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid http.Status: %s", resp.Status)
	}

	var result struct {
		util.Error
		Token
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(data))

	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if result.ErrCode != util.ErrCodeOK {
		return nil, result.Error
	}

	tk := new(Token)
	*tk = result.Token
	tk.CreatedAt = time.Now().Unix()

	c.Token = tk

	return tk, nil
}

// 不推荐使用
//func (c *Client) RefreshToken(code string) (*Token, error) {
//	return nil, nil
//}

func (c *Client) exchangeURL(code string) string {
	return fmt.Sprintf("%s?appid=%s"+
		"&secret=%s"+
		"&code=%s"+
		"&grant_type=authorization_code",
		c.Endpoint.TokenURL,
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.AppSecret),
		url.QueryEscape(code),
	)
}
