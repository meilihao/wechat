package oauth2

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/meilihao/wechat/oauth2"
)

type Config oauth2.Config

func (c *Config) AuthCodeURL(state string) string {
	return fmt.Sprintf("%s?appid=%s" +
		+"&redirect_uri=%s" +
		"&response_type=code" +
		"&scope=%s" +
		"&state=%s#wechat_redirect",
		c.Endpoint.AuthURL,
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.RedirectURL),
		url.QueryEscape(strings.Join(c.Scopes, ",")),
		url.QueryEscape(state),
	)
}
