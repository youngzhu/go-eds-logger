package logger_test

import (
	"goeds/logger"
	"testing"
)

func TestDailyReport(t *testing.T) {

	logger.RetrieveWorkReportRandom()

	err := logger.DailyReport("2025-10-12")
	if err != nil {
		t.Fatal(err)
	}
}
