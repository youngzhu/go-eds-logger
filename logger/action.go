package logger

import (
	"github.com/youngzhu/godate"
	"log"
)

type actionLogger struct{}

func init() {
	var action actionLogger
	register("action", action)
}

// Execute 每周一自动执行，填写周报和5天的日报
func (a actionLogger) Execute() {
	mon := godate.NewDate()
	logWholeWeek(mon)

	// 周末调休
	sat, _ := mon.AddDay(5)
	sun, _ := mon.AddDay(6)

	extraDays := RetrieveExtraDays()

	for _, d := range []string{sat.String(), sun.String()} {
		if _, ok := extraDays[d]; ok {
			log.Println("调休", d)
			workLog(d)
		}
	}

}
