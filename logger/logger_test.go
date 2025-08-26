package logger

import (
	"github.com/youngzhu/go-smail"
	"github.com/youngzhu/godate"
	"github.com/youngzhu/godate/chinese"
	"testing"
	"time"
)

//func TestRetrieveLogContent(t *testing.T) {
//	var c1 LogContent
//	var c2 *LogContent
//
//	c2, err := RetrieveLogContent()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	c1 = *c2
//
//	t.Logf("c1:%v", c1)
//	t.Log("c2:", c2)
//}

func TestGodate_IsWorkday(t *testing.T) {
	// 测试是否可以获取工作日
	today := godate.Today()
	t.Logf("Today: %s", today)

	for i := 0; i < 7; i++ {
		// 获取当前日期
		date, _ := today.AddDay(i)
		t.Logf("Date: %s, Is Workday: %v", date, chinese.IsWorkDayInChina(date))
	}

	day10_1 := godate.MustDate(2025, 10, 1)
	t.Logf("Date: %s, Is Workday: %v", day10_1, day10_1.IsWorkday())
	t.Logf("Date: %s, Is Workday: %v", day10_1, chinese.IsWorkDayInChina(day10_1))
	t.Logf("Date: %s, Is Offday: %v", day10_1, chinese.IsOffDayInChina(day10_1))
}

func TestSmail(t *testing.T) {
	smail.SendMail(time.Now().String()+"成功", "")
}
