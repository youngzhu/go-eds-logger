package logger

import (
	"github.com/youngzhu/godate"
)

type actionLogger struct{}

func init() {
	var action actionLogger
	register("action", action)
}

// Execute 每周一自动执行，填写周报和5天的日报
func (a actionLogger) Execute() {
	logWholeWeek(godate.NewDate())
}
