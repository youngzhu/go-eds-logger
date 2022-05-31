package main

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestX(t *testing.T) {
	// log.Println(os.Getwd()) // 可以查看当前的工作目录
	// secret, err := secret.RetrieveSecret()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(secret)

	// 获取当前时间
	now := time.Now()
	todayStr := now.Format("2006-01-02")

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
		logDay := now.Add(time.Hour * time.Duration(24*i))
		// log.Println(fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day()))
		// "2006-01-02 15:04:05"
		// 这算小彩蛋吗？还必须这个时间才行。。。
		// log.Println(now.Format("2006-01-02"))
		// log.Println(time.Now().Format("2006-01-02 15:04:05"))

		log.Println("日报", logDay.Format("2006-01-02"))
	}

	time.Sleep(5 * time.Second)

	log.Println("周报", todayStr)

	// date := time.Date(2021, 9, 18, 13, 14, 59, 0, time.UTC)

}
