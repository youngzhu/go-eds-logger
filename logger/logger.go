package logger

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/youngzhu/godate"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

type dayTime struct {
	startTime string
	endTime   string
}

var (
	am = dayTime{startTime: "10:00", endTime: "12:00"}
	pm = dayTime{startTime: "13:00", endTime: "18:00"}
)

func Daily(logUrl, projectID, logDate, logContent string) error {
	logUrl = logUrl + "&LogDate=" + logDate
	log.Println(logUrl)

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams, err := getHiddenParams(logUrl)
	if err != nil {
		return err
	}
	//fmt.Println(hiddenParams)

	for _, t := range []dayTime{am, pm} {
		err := doWorkLog(logUrl, projectID, logDate, logContent, t, hiddenParams)
		if err != nil {
			return fmt.Errorf("日志操作失败：%w", err)
		}
	}

	log.Println("日志操作成功", logDate)
	time.Sleep(800 * time.Millisecond)

	return nil
}

func doWorkLog(workLogUrl, projectID, logDate, logContent string, dt dayTime, hiddenParams map[string]string) error {
	startTime, endTime := dt.startTime, dt.endTime

	logParams := url.Values{}
	logParams.Set("__EVENTTARGET", "hplbWorkType")
	logParams.Set("__EVENTARGUMENT", "")
	logParams.Set("__LASTFOCUS", "")
	logParams.Set("__VIEWSTATEGENERATOR", "3A8BE513")
	logParams.Set("txtDate", logDate)
	logParams.Set("txtStartTime", startTime)
	logParams.Set("txtEndTime", endTime)
	logParams.Set("ddlProjectList", projectID)
	logParams.Set("hplbWorkType", "0106")
	logParams.Set("hplbAction", "010601")
	logParams.Set("TextBox1", "")
	logParams.Set("txtMemo", logContent)
	logParams.Set("btnSave", "+%E7%A1%AE+%E5%AE%9A+")
	logParams.Set("txtnodate", logDate)
	logParams.Set("txtnoStartTime", startTime)
	logParams.Set("txtnoEndTime", endTime)
	logParams.Set("TextBox6", "")
	logParams.Set("txtnoMemo", "")
	logParams.Set("txtCRMDate", logDate)
	logParams.Set("txtCRMStartTime", startTime)
	logParams.Set("txtCRMEndTime", endTime)
	logParams.Set("TextBox5", "")
	logParams.Set("txtCRMMemo", "")

	for key, value := range hiddenParams {
		logParams.Set(key, value)
	}

	//fmt.Println(logParams)
	_, err := DoPost(workLogUrl, strings.NewReader(logParams.Encode()))
	return err
}

func doWeeklyLog(logUrl, logDate string, lc LogContent) error {
	logParams := url.Values{}
	logParams.Set("hidCurrRole", "")
	logParams.Set("hidWeeklyState", "")
	logParams.Set("WeekReportDate", logDate)
	logParams.Set("txtWorkContent", lc.WeeklyWorkContent)
	logParams.Set("txtStudyContent", lc.WeeklyStudyContent)
	logParams.Set("txtSummary", lc.WeeklySummary)
	logParams.Set("txtPlanWork", lc.WeeklyPlanWork)
	logParams.Set("txtPlanStudy", lc.WeeklyPlanStudy)
	logParams.Set("btnSubmit", "%E6%8F%90%E4%BA%A4")

	// 周报没有校验
	// 通过get获取一些隐藏参数，用作后台校验
	//hiddenParams, err := getHiddenParams(logUrl)
	//if err != nil {
	//	return err
	//}
	//for key, value := range hiddenParams {
	//	logParams.Set(key, value)
	//}

	_, err := DoPost(logUrl, strings.NewReader(logParams.Encode()))
	if err != nil {
		return err
	}

	log.Println("周报填写成功", logDate)
	time.Sleep(2 * time.Second)

	return nil
}

func getHiddenParams(getUrl string) (map[string]string, error) {
	result := make(map[string]string)

	respHtml, err := DoGet(getUrl)
	if err != nil {
		//log.Println("getHiddenParams error:", err)
		return nil, fmt.Errorf("获取参数失败：%w", err)
	}
	//println(respHtml)

	keys := []string{"__EVENTVALIDATION", "__VIEWSTATE"}

	for _, k := range keys {
		result[k] = getValueFromHtml(respHtml, k)
	}

	return result, nil
}

func getValueFromHtml(html, key string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		log.Fatalln(err)
	}

	var value = ""
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		if id == key {
			value, _ = s.Attr("value")
			// fmt.Println("i", i, "选中的文本", value)
			return
		}

	})

	return value
}
func DoWeekly(urlWeekly, urlDaily, projectID string, lc LogContent) error {
	today := godate.Today()
	workdays := today.Workdays()

	// 先写周报
	// 只能填写本周周报（周一）!!!
	err := doWeeklyLog(urlWeekly, workdays[0].String(), lc)
	if err != nil {
		return err
	}

	// 再写日报
	for _, day := range workdays {
		err = Daily(urlDaily, day.String(), projectID, lc.DailyWorkContent)
		if err != nil {
			return err
		}
	}

	// 周末调休
	sat, _ := today.AddDay(5)
	sun, _ := today.AddDay(6)

	extraDays := RetrieveExtraDays()

	for _, dd := range []string{sat.String(), sun.String()} {
		if _, ok := extraDays[dd]; ok {
			log.Println("调休", dd)
			err = Daily(urlDaily, dd, projectID, lc.DailyWorkContent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// RetrieveExtraDays
// 返回map方便使用（查找）
func RetrieveExtraDays() map[string]struct{} {
	days := make(map[string]struct{})

	f, err := os.Open("data/extraDays.txt")
	defer f.Close()
	if err != nil {
		log.Println(err)
		return days
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		days[scanner.Text()] = struct{}{}
	}

	return days
}
