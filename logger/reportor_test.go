package logger_test

import (
	"goeds/logger"
	"testing"
)

func TestDailyReport(t *testing.T) {
	logger.DailyReport("2025-09-29")
}
