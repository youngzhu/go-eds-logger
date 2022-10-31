package main

import (
	"github.com/youngzhu/go-eds-logger/logger"
	"github.com/youngzhu/godate"
	"log"
)

func main() {

	err := logger.Login()
	if err != nil {
		sendFailedMail(err.Error())
		log.Fatalln(err)
		// os.Exit(1)
	}

	err = logger.PrepareData()
	if err != nil {
		sendFailedMail(err.Error())
		log.Fatalln(err)
	}

	//logger.Run()

	sendSuccessfulMail()
}

var today = godate.Today()

func sendSuccessfulMail() {
	logger.SendMail(today.String()+"成功", "RT")
}

func sendFailedMail(errMsg string) {
	logger.SendMail(today.String()+"失败", errMsg)
}
