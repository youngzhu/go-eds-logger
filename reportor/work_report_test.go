package reportor

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/youngzhu/godate"
	"log"
	"testing"
)

func Test_retrieveWorkReportFromInternet(t *testing.T) {
	workReport, err := retrieveWorkReportFromInternet()

	if err != nil {
		t.Error(err)
	}

	if workReport.LastWeekWorkContent == "" {
		t.Error("LastWeekWorkContent should not empty.")
	}
}

func Test_retrieveWorkReportFromLocal(t *testing.T) {
	workReport := retrieveWorkReportFromLocal()

	if workReport.LastWeekWorkContent == "" {
		t.Error("LastWeekWorkContent should not empty.")
	}
}

func TestWorkReport_formatDailyReportUrl(t *testing.T) {
	reportUrl := fmt.Sprintf("%s&LogDate=%s", viper.GetString("urls.daily"), "2023-10-01")

	want := "http://eds.newtouch.cn/eds3/worklog.aspx?tabid=0&LogDate=2023-10-01"
	if reportUrl != want {
		t.Errorf("Expected URL to be '%s', got '%s'", want, reportUrl)
	}
}

func TestWorkReport_getHiddenParams(t *testing.T) {
	// 先登录
	err := login()
	if err != nil {
		t.Error("登录失败：", err)
	}

	reportUrl := fmt.Sprintf("%s&LogDate=%s", viper.GetString("urls.daily"), godate.Today())
	log.Println(reportUrl)

	_, err = getHiddenParams(reportUrl)
	if err != nil {
		t.Error("getHiddenParams error:", err)
	}

}
