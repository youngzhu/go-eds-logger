package logger_test

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/youngzhu/godate"
	"github.com/youngzhu/godate/chinese"
	"goeds/logger"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 在所有测试运行前执行一次（类似 @BeforeClass）
	fmt.Println("Global setup - runs once before all tests")

	// 执行初始化操作
	//setupDatabase()
	//loadConfig()
	logger.InitConfig()

	// 运行所有测试
	exitCode := m.Run()

	// 在所有测试运行后执行（类似 @AfterClass）
	//teardownDatabase()

	os.Exit(exitCode)
}

func TestDailyReport(t *testing.T) {

	err := logger.LoginX()
	if err != nil {
		t.Fatal(err)
	}

	logger.LoadWorkReportRandomly()

	err = logger.DailyReport("2026-05-11")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoginX(t *testing.T) {
	err := logger.LoginX()
	if err != nil {
		t.Fatal(err)
	}
}

func TestInitConfig(t *testing.T) {
	err := logger.InitConfig()
	if err != nil {
		t.Fatal(err)
	}

	testVal := viper.GetString("test")
	if testVal != "123" {
		t.Fatalf("Expected '123', got '%s'", testVal)
	}

	home := viper.GetString("urls.home")
	if home == "" {
		t.Fatal("Expected non-empty 'urls.home'")
	}

}

func TestReportor_QueryDailyReport(t *testing.T) {
	logger.LoginX()
	//logger.QueryDailyReport("2025-10-27")
	report, _ := logger.QueryDailyReport("2026-05-11")
	if (report == logger.DailyReportDetail{}) {
		t.Fatal("Expected non-empty report, got empty")
	}
}

func TestReportor_HasReport(t *testing.T) {
	logger.LoginX()

	var testcases = []struct {
		reportDate string
		except     bool
	}{
		{"2025-11-09", false},
		{"2025-11-10", true},
	}

	for _, testcase := range testcases {
		got := logger.HasReport(testcase.reportDate)
		if got != testcase.except {
			t.Errorf("got: %v, want: %v", got, testcase.except)
		}
	}

}

func TestAWeek(t *testing.T) {
	today := godate.Today()
	monday := chinese.NewCNDate(today.Workdays()[0])

	date := monday
	for {
		log.Println(date, "-", date.Name)

		date = date.NextDay()

		if date.Weekday() == godate.Monday {
			break
		}
	}
}
