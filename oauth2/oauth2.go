package oauth2

import (
	"time"

	"github.com/meilihao/wechat/util"
)

type Endpoint struct {
	AuthURL         string
	TokenURL        string
	RefreshTokenURL string
}

// from golang.org/x/oauth2
type Config struct {
	// AppID is the application's ID.
	AppID string

	// AppSecret is the application's secret.
	AppSecret string

	// Endpoint contains the resource server's token endpoint
	// URLs.
	Endpoint Endpoint

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string

	// Scope specifies optional requested permissions.
	Scopes []string
}

/*
{
    "access_token":"xxx",
    "expires_in":7200,
    "refresh_token":"xxx",
    "openid":"xxx",
    "scope":"snsapi_base"
	"unionId":"xxx"
}
*/
type Token struct {
	AccessToken  string `json:"access_token"`            // 授权接口调用凭证, 注意：此access_token与基础支持的access_token不同
	CreatedAt    int64  `json:"-"`                       // 获取到token的时间
	ExpiresIn    int64  `json:"expires_in"`              // access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token,omitempty"` // 用户刷新access_token

	OpenId  string `json:"openid,omitempty"`  // 用户唯一标识，请注意，在未关注或关注公众号时，用户访问公众号的授权网页后都会产生同一个OpenID
	UnionId string `json:"unionid,omitempty"` // 当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段
	Scope   string `json:"scope,omitempty"`   // 用户授权的作用域，使用逗号（,）分隔
}

// 判断 token.AccessToken 是否已过期
// 7200 * 5% = 360
func (token *Token) IsExpired() bool {
	return time.Now().Unix() >= token.CreatedAt+token.ExpiresIn-util.TimeBufferForAccessToken
}

func (token *Token) IsValidRefreshToken() bool {
	return token.CreatedAt > 0 && time.Now().Unix() <= token.CreatedAt+(30*24)*3600-util.TimeBufferForRefreshToken
}
