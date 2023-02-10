package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrConnection  = errors.New("连接错误")
	ErrInvalidUser = errors.New("用户名或密码错误")
)

func newClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func CheckURL(url string) error {
	resp, err := newClient().Head(url) // 只请求网站的 http header信息
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w, statusCode: %d", ErrConnection, resp.StatusCode)
	}
	return nil
}

func Login(loginURL, userId, password string) error {
	// data := `{"UserId":"###", "UserPsd":"***"}`
	// data := "UserId=###&UserPsd=***"
	// params := url.Values{
	// 	"UserId":  {"###"},
	// 	"UserPsd": {"***"},
	// }
	params := url.Values{}
	params.Add("UserId", userId)
	params.Add("UserPsd", password)
	// var request *http.Request
	// request, err = http.NewRequest(http.MethodPost, URL_LOGIN, strings.NewReader(data))
	// request, err = http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(params.Encode()))

	resp, err := doRequest(loginURL, http.MethodPost,
		strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	if strings.Contains(resp, ErrInvalidUser.Error()) {
		return ErrInvalidUser
	}

	log.Println("登陆成功")
	
	return nil
}
