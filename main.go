package main

import (
	"github.com/youngzhu/go-eds-logger/logger"
	"log"
)

func main() {
	err := logger.Login()
	if err != nil {
		// 正常返回还不行，需要有错误发送邮件通知
		// return
		log.Fatalln(err)
		// os.Exit(1)
	}

	logger.Run()
}
