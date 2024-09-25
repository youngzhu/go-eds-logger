package main

import (
	"edser/config"
	"edser/logger"
	smail "github.com/youngzhu/go-smail"
	"github.com/youngzhu/godate"
	"log"
	"time"
)

var cfg config.Configuration

func init() {
	var err error
	c, err := config.Load("config.json")
	if err != nil {
		panic(err)
	}
	cfg = c
}

func main() {
	// 如果有一个失败，全部重头再来
	var err error
	// github action 最多执行6分钟，所以尝试5次
	for i := 0; i < 5; i++ {
		err = logger.Run(cfg)
		if err == nil {
			// 成功
			break
		} else {
			log.Printf("第 %d 次失败: %s", i+1, err.Error())
			time.Sleep(time.Minute)
		}
	}

	if err != nil {
		sendFailedMail(err.Error())
		log.Fatalln(err) // 结束
	}

	sendSuccessfulMail()
}

var today = godate.Today()

func sendSuccessfulMail() {
	smail.SendMail(today.String()+"成功", "")
}

func sendFailedMail(errMsg string) {
	smail.SendMail(today.String()+"失败", errMsg)
}
