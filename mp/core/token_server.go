package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/meilihao/wechat/util"
)

var ErrNoToken = errors.New("no base token")

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	CreatedAt   int64  `json:"-"`
}

func (token *Token) IsExpired() bool {
	return time.Now().Unix() >= token.CreatedAt+token.ExpiresIn-6*60
}

// 其他业务逻辑服务器所使用的access_token by AppID
type TokenStorage interface {
	Token(string) (*Token, error)
	SetToken(string, *Token) error
}

type TokenServer struct {
	AppID      string
	AppSecret  string
	HttpClient *http.Client

	TokenStorage
	token *Token

	timer            *time.Timer
	lock             *sync.RWMutex
	refreshTokenChan chan struct{}
}

func NewTokenSever(appID, appSecret string, httpClient *http.Client, store TokenStorage) *TokenServer {
	ts := &TokenServer{
		AppID:      appID,
		AppSecret:  appSecret,
		HttpClient: http.DefaultClient,
		lock:       new(sync.RWMutex),
	}

	if httpClient != nil {
		ts.HttpClient = httpClient
	}

	if store != nil {
		ts.TokenStorage = store
	}

	ts.timer = time.NewTimer(getTimerDuration(7200)) // set timer, Avoid select{} use ts.timer at first
	ts.refreshTokenChan = make(chan struct{}, 1)
	ts.refreshTokenChan <- struct{}{} // first updateToken

	go ts.updateToken()

	return ts
}

func (ts *TokenServer) InitToken() (int64, error) {
	resp, err := ts.HttpClient.Get(fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		ts.AppID,
		ts.AppSecret,
	))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Invalid http.Status: %s", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	log.Println(string(data))

	var result struct {
		util.Error
		Token
	}

	if err = json.Unmarshal(data, &result); err != nil {
		return 0, err
	}

	if result.ErrCode != util.ErrCodeOK {
		return 0, result.Error
	}

	ts.lock.Lock()
	defer ts.lock.Unlock()

	tk := new(Token)
	*tk = result.Token
	tk.CreatedAt = time.Now().Unix()
	ts.token = tk

	if ts.TokenStorage != nil {
		if err = ts.TokenStorage.SetToken(ts.AppID, tk); err != nil {
			log.Printf("set token err: %v\n", err)
		}
	}

	return tk.ExpiresIn, nil
}

func (ts *TokenServer) updateToken() {
	var du int64
	var err error

	for {
		select {
		case <-ts.refreshTokenChan:
			du, err = ts.InitToken()
		case <-ts.timer.C:
			du, err = ts.InitToken()
		}

		if du > 0 {
			log.Printf("reset timer(%d,%v): %t", du, time.Now(), ts.timer.Reset(getTimerDuration(du)))
		} else {
			log.Printf("err: %v\n", err)
			time.Sleep(1000 * time.Millisecond)

			ts.refreshTokenChan <- struct{}{}
		}
	}
}

func (ts *TokenServer) Token() *Token {
	ts.lock.RLock()

	var tk *Token
	if ts.token != nil && !ts.token.IsExpired() {
		tk = new(Token)
		*tk = *ts.token
	}

	ts.lock.RUnlock()

	return tk
}

func (ts *TokenServer) AccessToken() string {
	ts.lock.RLock()

	var at string
	if ts.token != nil && !ts.token.IsExpired() {
		at = ts.token.AccessToken
	}

	ts.lock.RUnlock()

	return at
}

func getTimerDuration(expiresIn int64) time.Duration {
	return time.Second * time.Duration(expiresIn-util.TimeBufferForAccessToken)
}
