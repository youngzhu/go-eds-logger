package logger

import (
	"time"
)

type actionLogger struct{}

func init() {
	var action actionLogger
	register("action", action)
}

// Execute 每周一执行，填写周报和5天的日报
func (a actionLogger) Execute() {
	logDate := time.Now()

	var workday []string
	for i := 0; i < 5; i++ {
		workday = append(workday, logDate.Format("2006-01-02"))
		logDate = logDate.Add(time.Hour * 24)
	}

	// 先写周报
	// 只能填写本周周报（周一）!!!
	workWeeklyLog(workday[0])

	time.Sleep(5 * time.Second)

	// 再写日报
	for _, day := range workday {
		workLog(day)
		time.Sleep(1 * time.Second)
	}
}
