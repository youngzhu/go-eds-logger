package logger_test

import (
	"goeds/logger"
	"testing"
)

func TestDailyReport(t *testing.T) {

	err := logger.LoginX("12234", "young12234")
	if err != nil {
		t.Fatal(err)
	}

	logger.RetrieveWorkReportRandom()

	err = logger.DailyReport("2025-10-13")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoginX(t *testing.T) {
	err := logger.LoginX("12234", "young12234")
	if err != nil {
		t.Fatal(err)
	}
}
