package main

import (
	"github.com/youngzhu/godate"
	"goeds/config"
	"goeds/logger"
	"log"
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
	err := logger.Run(cfg)
	if err != nil {
		sendFailedMail(err.Error())
		log.Fatalln(err) // 结束
	}

	sendSuccessfulMail()
}

var today = godate.Today()

func sendSuccessfulMail() {
	logger.SendMail(today.String()+"成功", "RT")
}

func sendFailedMail(errMsg string) {
	logger.SendMail(today.String()+"失败", errMsg)
}
