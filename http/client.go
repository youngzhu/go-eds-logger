package http

import (
	"errors"
	"fmt"
	"net/http"
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

func Connect(url string) error {
	resp, err := newClient().Head(url) // 只请求网站的 http header信息
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w, statusCode: %d", ErrConnection, resp.StatusCode)
	}
	return nil
}
