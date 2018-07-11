package oauth2

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/meilihao/wechat/oauth2"
)

type Config oauth2.Config

func (c *Config) AuthCodeURL(state string) string {
	return fmt.Sprintf("%s?appid=%s"+
		"&redirect_uri=%s"+
		"&response_type=code&scope=%s"+
		"&state=%s"+
		"&connect_redirect=1#wechat_redirect",
		c.Endpoint.AuthURL,
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.RedirectURL),
		url.QueryEscape(strings.Join(c.Scopes, ",")),
		url.QueryEscape(state),
	)
}

func (c *Config) AuthCodeURLWithCustom(state string, us url.Values) string {
	tmp := c.RedirectURL
	if len(us) > 0 {
		tmp += "?" + us.Encode()
	}

	return fmt.Sprintf("%s?appid=%s"+
		"&redirect_uri=%s"+
		"&response_type=code&scope=%s"+
		"&state=%s"+
		"&connect_redirect=1#wechat_redirect",
		c.Endpoint.AuthURL,
		url.QueryEscape(c.AppID),
		url.QueryEscape(tmp),
		url.QueryEscape(strings.Join(c.Scopes, ",")),
		url.QueryEscape(state),
	)
}
