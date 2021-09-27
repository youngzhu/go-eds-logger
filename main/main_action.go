package main

import (
	"log"
	"time"
	"flag"

	"github.com/youngzhu/go-eds-logger/logger"
	"github.com/youngzhu/go-eds-logger/secret"
)

var (
	userId string
	userPsd string
	cookie string
)

func main() {
	
	flag.StringVar(&userId, "i", "", "login user id")
	flag.StringVar(&userPsd, "p", "", "login password")
	flag.StringVar(&cookie, "c", "", "cookie")

	// 这句不能省，
	// 还以为用上面的方式已经将参数赋给了变量
	flag.Parse() // 解析入参

	// log.Println("id", userId)
	// log.Println("pwd", userPsd)
	// log.Println("cookie", cookie)

	loginInfo := secret.Secret{userId, userPsd, cookie}
	err := logger.Login(&loginInfo)
	if err != nil {
		// 正常返回还不行，需要有错误发送邮件通知
		// return
		log.Fatalln("网站服务错误", err)
		// os.Exit(1)
	}

	logger.LogFromSpecifiedDay(time.Now())
}