package logger

import (
	"log"
	"time"
)

type manualLogger struct{}

func init() {
	var manual manualLogger
	register("manual", manual)
}

// Execute 手动执行指定日期的日志
func (a manualLogger) Execute() {
	logFrom, _ := time.Parse("2006-01-02", "2022-05-16")
	log.Println(logFrom)

	// 填一周
	logWholeWeek(logFrom)

	// 填一天
	logTheSpecifiedDay(logFrom)
	// LogTheSpecifiedDay(time.Now())
}
