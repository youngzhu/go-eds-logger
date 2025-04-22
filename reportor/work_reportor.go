package reportor

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strings"
)

// 学 viper 设置一个影子变量
var _reportor *WorkReportor

func init() {
	_reportor = New()
}

type WorkReportor struct {
	//urls map[string]string
	workReport WorkReport // 周报内容
}

func New() *WorkReportor {
	wr := new(WorkReportor)

	//wr.urls = make(map[string]string)

	return wr
}

// 登录
func login() (err error) {
	return _reportor.login()
}

func (r *WorkReportor) login() (err error) {
	userID := viper.GetString("usr-id")
	passcode := viper.GetString("usr-pwd")
	if userID == "" || passcode == "" {
		return errors.New("用户名或密码不能为空")
	}

	params := url.Values{}
	params.Set("UserId", userID)
	params.Set("UserPsd", passcode)

	resp, err := r.doPost(viper.GetString("urls.login"), strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("登录错误：%w", err)
	}

	if strings.Contains(resp, ErrInvalidUser.Error()) {
		return ErrInvalidUser
	}

	log.Println("登陆成功")

	return
}

// 装载周报内容
func loadWorkReport() {
	_reportor.loadWorkReport()
}

func (r *WorkReportor) loadWorkReport() {
	r.workReport = retrieveWorkReport()
}
