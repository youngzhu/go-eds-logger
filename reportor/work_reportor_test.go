package reportor

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/youngzhu/godate"
	"log"
	"testing"
)

func Test_login(t *testing.T) {
	err := login()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_formatDailyReportUrl(t *testing.T) {
	reportUrl := fmt.Sprintf("%s&LogDate=%s", viper.GetString("urls.daily"), "2023-10-01")

	want := "http://eds.newtouch.cn/eds3/worklog.aspx?tabid=0&LogDate=2023-10-01"
	if reportUrl != want {
		t.Errorf("Expected URL to be '%s', got '%s'", want, reportUrl)
	}
}

func Test_getHiddenParams(t *testing.T) {
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
