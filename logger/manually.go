package logger

import (
	"edser/config"
	"github.com/youngzhu/godate"
)

// 尝试过为了简化docker的引入，将它放入单独的目录下
// 结果发生编译错误

type manualLogger struct{}

func init() {
	var manual manualLogger
	register("manual", manual)
}

// Execute 手动执行指定日期的日志
func (a manualLogger) Execute(cfg config.Configuration) (err error) {
	// 填一周
	err = logWholeWeek(cfg, godate.Today())

	// 填一天
	//logDay, _ := godate.NewDateYMD(2022, 6, 6)
	//logDay, _ := godate.Today().AddDay(0)
	//log.Println(logDay)
	//workLog(cfg, logDay.String())

	return
}
