package main

import (
	"github.com/youngzhu/go-eds-logger/logger"
	"github.com/youngzhu/go-eds-logger/secret"
	"log"
	"os"
)

func main() {
	loginInfo := secret.Secret{UserId: os.Getenv("USER_ID"), UserPsd: os.Getenv("USER_PWD")}
	err := logger.Login(&loginInfo)
	if err != nil {
		// 正常返回还不行，需要有错误发送邮件通知
		// return
		log.Fatalln("网站服务错误", err)
		// os.Exit(1)
	}

	//logger.LogFromSpecifiedDay(time.Now())
}
