package logger

import (
	"github.com/youngzhu/godate"
	"log"
)

// 尝试过为了简化docker的引入，将它放入单独的目录下
// 结果发生编译错误

type manualLogger struct{}

func init() {
	var manual manualLogger
	register("manual", manual)
}

// Execute 手动执行指定日期的日志
func (a manualLogger) Execute() {
	//logFrom, _ := godate.NewDateYMD(2022, 6, 6)
	logFrom := godate.Today()
	log.Println(logFrom)

	// 填一周
	logWholeWeek(logFrom)

	// 填一天
	//logTheSpecifiedDay(logFrom)
	// LogTheSpecifiedDay(time.Now())
}
