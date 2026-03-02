/*
Copyright © 2023 youngzy
Copyrights apply to this source code.
Check LICENSE for details.

*/
package main

import (
	"github.com/youngzhu/go-smail"
	"github.com/youngzhu/godate"
	"goeds/cmd"
	"log"
	"math/rand"
	"time"
)

func init() {
	// 在 init 函数中设置一次种子（程序启动时执行）
	rand.Seed(time.Now().UnixNano())
}

func main() {
	err := cmd.Execute()

	if err != nil {
		sendFailedMail(err.Error())
		log.Fatalln(err) // 结束
	}

	sendSuccessfulMail()
}

var today = godate.Today()

func sendSuccessfulMail() {
	smail.SendMail(today.String()+"任务计划成功", "")
}

func sendFailedMail(errMsg string) {
	smail.SendMail(today.String()+"任务计划失败", errMsg)
}
