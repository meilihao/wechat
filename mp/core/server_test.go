package core

import (
	"net/url"
	"testing"

	"github.com/meilihao/wechat/util"
)

func TestDecryptMsg(t *testing.T) {
	raw := []byte(`<xml>
	                <ToUserName><![CDATA[wx5823bf96d3bd56c7]]></ToUserName>
	                <Encrypt><![CDATA[RypEvHKD8QQKFhvQ6QleEB4J58tiPdvo+rtK1I9qca6aM/wvqnLSV5zEPeusUiX5L5X/0lWfrf0QADHHhGd3QczcdCUpj911L3vg3W/sYYvuJTs3TUUkSUXxaccAS0qhxchrRYt66wiSpGLYL42aM6A8dTT+6k4aSknmPj48kzJs8qLjvd4Xgpue06DOdnLxAUHzM6+kDZ+HMZfJYuR+LtwGc2hgf5gsijff0ekUNXZiqATP7PF5mZxZ3Izoun1s4zG4LUMnvw2r+KqCKIw+3IQH03v+BCA9nMELNqbSf6tiWSrXJB3LAVGUcallcrw8V2t9EL4EhzJWrQUax5wLVMNS0+rUPA3k22Ncx4XXZS9o0MBH27Bo6BpNelZpS+/uh9KsNlY6bHCmJU9p8g7m3fVKn28H3KDYA5Pl/T8Z1ptDAVe0lXdQ2YoyyH2uyPIGHBZZIs2pDBS8R07+qN+E7Q==]]></Encrypt>
                </xml>`)
	aesKey := "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C"
	token := "QDG6eK"
	appId := "wx5823bf96d3bd56c7"

	us := url.Values{}
	us.Set("signature", "477715d11cdb4164915debcba66cb864d751f3e6")
	us.Set("timestamp", "1409659813")
	us.Set("nonce", "1372623149")
	data, err := DecryptMsg(raw, util.DecodeAESKey(aesKey), appId, token, us)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("data: %s\n", string(data))
	}
}
