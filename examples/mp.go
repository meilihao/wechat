package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/sid"
	"github.com/gin-gonic/gin"
	"github.com/meilihao/wechat/mp/core"
	"github.com/meilihao/wechat/mp/message/template"
	mpoauth2 "github.com/meilihao/wechat/mp/oauth2"
	"github.com/meilihao/wechat/oauth2"
	"github.com/meilihao/wechat/util"
	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger

	wxAppId        = "wxxxx" // 测试号的appID
	wxAppSecret    = "xxx"
	Token          = "xxxServerSettingToken" // 服务器配置里的token
	TestAppToken   = "xxxTestToken"          // 测试号使用的Token
	TestAppURL     = "/"                     // 测试号使用的URL
	maxAge         = 24 * 3600 * 365 * 100
	sessionStorage = session.New(20*60, 60*60)
	mpOauth2Config = &mpoauth2.Config{
		AppID:       wxAppId,
		AppSecret:   wxAppSecret,
		RedirectURL: "http://opengolang.com/callback/wechat/mp/oauth2",
		Scopes:      []string{"snsapi_base"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://open.weixin.qq.com/connect/oauth2/authorize",
			TokenURL: "https://api.weixin.qq.com/sns/oauth2/access_token",
		},
	}
	tokenServer = core.NewTokenSever(wxAppId, wxAppSecret, nil, nil)
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar = logger.Sugar()

	r := gin.Default()

	r.POST(TestAppURL, _Callback)
	r.GET("/callback/wechat/mp", _VerifySetting)

	r.GET("/redirect/wechat/mp/oauth2", _AuthFromWeb)
	r.GET("/callback/wechat/mp/oauth2", _AuthFromWebCallback)

	r.Run(":8888")
}

func ErrJSON(err error) gin.H {
	return gin.H{
		"errmsg": err.Error(),
	}
}

func bodyReader(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, errors.New("EmptyBody")
	}
	defer req.Body.Close()

	return ioutil.ReadAll(req.Body)
}

// 接受微信服务器的事件推送
func _Callback(ctx *gin.Context) {
	content, err := bodyReader(ctx.Request)
	sugar.Info(
		err,
		util.CheckSignature(ctx.Query("signature"), TestAppToken, ctx.Query("timestamp"), ctx.Query("nonce")),
		string(content))
}

// 验证服务器配置, 请使用正式的服务器配置Token
func _VerifySetting(ctx *gin.Context) {
	// checkSignature
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	signature := ctx.Query("signature")
	echostr := ctx.Query("echostr")

	if util.CheckSignature(signature, Token, timestamp, nonce) {
		ctx.String(http.StatusOK, echostr)
		sugar.Debug("match signature")
	} else {
		sugar.Debug("no match signature")
	}
}

// 获取网页授权的RedirectURL
func _AuthFromWeb(ctx *gin.Context) {
	sessionID, _ := ctx.Cookie("sid")
	if sessionID == "" {
		sessionID = sid.New()

		ctx.SetCookie("sid", sessionID, maxAge, "/", "", false, true)
	}
	state := string(rand.NewHex())

	sugar.Debug("sid:", sessionID)
	sugar.Debug("state:", state)

	if err := sessionStorage.Set(sessionID, state); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)

		sugar.Error(err)
		return
	}

	authCodeURL := mpOauth2Config.AuthCodeURLWithCustom(state, nil)
	sugar.Debug("AuthCodeURL:", authCodeURL)

	ctx.Redirect(http.StatusFound, authCodeURL)
}

// 处理网页授权的回调
func _AuthFromWebCallback(ctx *gin.Context) {
	sugar.Debug(ctx.Request.URL.Query())

	sessionID, err := ctx.Cookie("sid")
	if sessionID == "" {
		sugar.Error(err)

		ctx.JSON(http.StatusBadRequest, ErrJSON(errors.New("No sid")))
		return
	}

	session, err := sessionStorage.Get(sessionID)
	if err != nil {
		sugar.Error(err)

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	savedState := session.(string)

	queryState := ctx.Query("state")
	if queryState == "" {
		ctx.JSON(http.StatusBadRequest, ErrJSON(errors.New("state 参数为空")))
		return
	}
	if savedState != queryState {
		ctx.JSON(http.StatusBadRequest, ErrJSON(fmt.Errorf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, queryState)))
		return
	}
	sessionStorage.Delete(sessionID)

	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, ErrJSON(errors.New("用户禁止授权")))
		return
	}

	oauth2Client := oauth2.NewClient(mpOauth2Config.AppID, mpOauth2Config.AppSecret,
		mpOauth2Config.Endpoint, nil)
	token, err := oauth2Client.Exchange(code)
	if err != nil {
		sugar.Error(err)

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	sugar.Debugf("token: %+v", token)

	if _, err = template.Send(tokenServer.AccessToken(), &template.TemplateMessage{
		ToUser:     token.OpenId,
		TemplateId: "B0pRhyJ_4skw5qFB-L6wG47R4Z4oQk7rfuWmWy2q95c",
		Data: json.RawMessage([]byte(`{
                   "first": {
                       "value":"恭喜你购买成功！",
                       "color":"#173177"
                   },
                   "keynote1":{
                       "value":"巧克力",
                       "color":"#173177"
                   },
                   "keynote2": {
                       "value":"39.8元",
                       "color":"#173177"
                   },
                   "keynote3": {
                       "value":"2014年9月22日",
                       "color":"#173177"
                   },
                   "remark":{
                       "value":"欢迎再次购买！",
                       "color":"#173177"
                   }
           }`)),
	}); err != nil {
		sugar.Error(err)

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.String(http.StatusOK, token.OpenId)
}
