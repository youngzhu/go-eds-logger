package logger

import (
	"edser/config"
	"github.com/youngzhu/godate"
)

type actionLogger struct{}

func init() {
	var action actionLogger
	register("Action", action)
}

// Execute 每周一自动执行，填写周报和5天的日报
func (a actionLogger) Execute(cfg config.Configuration) {
	mon := godate.NewDate()
	logWholeWeek(cfg, mon)
}
