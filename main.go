package main

import (
	"edser/config"
	"edser/logger"
	smail "github.com/youngzhu/go-smail"
	"github.com/youngzhu/godate"
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
	smail.SendMail(today.String()+"成功", "")
}

func sendFailedMail(errMsg string) {
	smail.SendMail(today.String()+"失败", errMsg)
}
