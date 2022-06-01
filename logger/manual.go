package logger

type manualLogger struct{}

func init() {
	var manual manualLogger
	register("manual", manual)
}

// Execute 手动执行指定日期的日志
func (a manualLogger) Execute() {

}
