package logger

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/youngzhu/godate"
	"goeds/http"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

type user struct {
	id       string
	password string
}

func getUser() (*user, error) {
	var loginID, loginPassword string
	var err error

	if loginID, err = getSecret("EDS_USR_ID"); err != nil {
		return nil, err
	}
	if loginPassword, err = getSecret("EDS_USR_PWD"); err != nil {
		return nil, err
	}

	return &user{id: loginID, password: loginPassword}, nil
}

const secretErr = "变量[%s]未配置\n"

func getSecret(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf(secretErr, key)
		return "", fmt.Errorf(secretErr, key)
	}
	return val, nil
}

var (
	ddlProjectList string
)

type dayTime struct {
	startTime string
	endTime   string
}

var am = dayTime{startTime: "10:00", endTime: "12:00"}
var pm = dayTime{startTime: "13:00", endTime: "18:00"}

func Daily(url, logDate, logContent string) error {
	url = url + "&LogDate=" + logDate

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams, err := getHiddenParams(url)
	if err != nil {
		return err
	}
	// fmt.Println(hiddenParams)

	for _, t := range []dayTime{am, pm} {
		err := doWorkLog(url, logDate, logContent, t, hiddenParams)
		if err != nil {
			return err
		}
	}

	log.Println("日志操作成功", logDate)
	time.Sleep(800 * time.Millisecond)

	return nil
}

func doWorkLog(workLogUrl, logDate, logContent string, dt dayTime, hiddenParams map[string]string) error {
	startTime, endTime := dt.startTime, dt.endTime

	logParams := url.Values{}
	logParams.Set("__EVENTTARGET", "hplbWorkType")
	logParams.Set("__EVENTARGUMENT", "")
	logParams.Set("__LASTFOCUS", "")
	logParams.Set("__VIEWSTATEGENERATOR", "3A8BE513")
	logParams.Set("txtDate", logDate)
	logParams.Set("txtStartTime", startTime)
	logParams.Set("txtEndTime", endTime)
	logParams.Set("ddlProjectList", ddlProjectList)
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

	//http.DoPost(workLogUrl, strings.NewReader(logParams.Encode()))
	return nil
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

	// 通过get获取一些隐藏参数，用作后台校验
	hiddenParams, err := getHiddenParams(logUrl)
	if err != nil {
		return err
	}
	for key, value := range hiddenParams {
		logParams.Set(key, value)
	}

	//http.DoPost(logUrl, strings.NewReader(logParams.Encode()))

	log.Println("周报填写成功", logDate)
	time.Sleep(2 * time.Second)

	return nil
}

func getHiddenParams(url string) (map[string]string, error) {
	result := make(map[string]string)

	respHtml, err := http.DoGet(url)
	if err != nil {
		//log.Println("getHiddenParams error:", err)
		return nil, err
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
func DoWeekly(urlWeekly, urlDaily string, lc LogContent) error {
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
		err = Daily(urlDaily, day.String(), lc.DailyWorkContent)
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
			err = Daily(urlDaily, dd, lc.DailyWorkContent)
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
