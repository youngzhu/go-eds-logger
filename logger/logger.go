package logger

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/youngzhu/godate"
	"github.com/youngzhu/godate/chinese"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

// 学 viper 设置一个影子变量
var lg *EDSLogger

func init() {
	lg = New()
}

type EDSLogger struct {
	projectID string // 项目编号
	urls      map[string]string
	cookie    string
	host      string

	lc         LogContent
	workReport WorkReport
}

func New() *EDSLogger {
	edsLogger := new(EDSLogger)

	edsLogger.urls = make(map[string]string)

	return edsLogger
}

func AddUrl(key, val string) {
	lg.AddUrl(key, val)
}
func (e *EDSLogger) AddUrl(key, val string) {
	e.urls[key] = val
}

func SetCookie(cookie string) {
	lg.SetCookie(cookie)
}
func (e *EDSLogger) SetCookie(cookie string) {
	e.cookie = cookie
}

func SetHost(host string) {
	lg.SetHost(host)
}
func (e *EDSLogger) SetHost(host string) {
	e.host = host
}

func Login(userId, password string) error {
	return lg.Login(userId, password)
}

func (e EDSLogger) Login(userId, password string) error {
	params := url.Values{}
	params.Set("UserId", userId)
	params.Set("UserPsd", password)

	resp, err := e.doPost(e.urls["login"], strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("登录错误：%w", err)
	}

	if strings.Contains(resp, ErrInvalidUser.Error()) {
		return ErrInvalidUser
	}

	log.Println("登陆成功")

	return nil
}

func RetrieveProjectID() error {
	return lg.RetrieveProjectID()
}
func (e *EDSLogger) RetrieveProjectID() error {
	respHtml, _ := e.doGet(e.urls["daily"])
	//println(respHtml)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(respHtml))

	if err != nil {
		return err
	}

	//fmt.Println("cookie:", e.cookie)

	var projectId string
	doc.Find("select").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		if id == "ddlProjectList" {
			projectId, _ = s.Children().Attr("value")
			return
		}
	})

	if projectId == "" {
		return errors.New("未能获取项目编号")
	}

	e.projectID = projectId

	return nil
}

//func ProjectID() string {
//	return lg.projectID
//}
//func (e *EDSLogger) ProjectID() string {
//	return e.projectID
//}

////
type dayTime struct {
	startTime string
	endTime   string
}

var (
	am = dayTime{startTime: "10:00", endTime: "12:00"}
	pm = dayTime{startTime: "13:00", endTime: "18:00"}
)

func DailyLog(logDate string) error {
	return lg.DailyLog(logDate)
}
func (e EDSLogger) DailyLog(logDate string) error {
	logUrl := e.urls["daily"]
	if logUrl == "" {
		return errors.New("logUrl为空")
	}

	logUrl = logUrl + "&LogDate=" + logDate

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams, err := getHiddenParams(logUrl)
	if err != nil {
		return err
	}
	//fmt.Println(hiddenParams)

	//log.Println("logContent:", e.lc.DailyWorkContent)

	workPlanDaily := e.workReport.workPlanDaily()

	for _, t := range []dayTime{am, pm} {
		//err := e.doWorkLog(logUrl, logDate, t, hiddenParams)
		err := e.doWorkLogSameDaily(logUrl, logDate, workPlanDaily, t, hiddenParams)
		if err != nil {
			return fmt.Errorf("日志操作失败：%w", err)
		}
	}

	log.Println("日志操作成功", logDate)
	time.Sleep(800 * time.Millisecond)

	return nil
}

func (e EDSLogger) doWorkLog(workLogUrl, logDate string, dt dayTime, hiddenParams map[string]string) error {
	startTime, endTime := dt.startTime, dt.endTime

	logParams := url.Values{}
	logParams.Set("__EVENTTARGET", "hplbWorkType")
	logParams.Set("__EVENTARGUMENT", "")
	logParams.Set("__LASTFOCUS", "")
	logParams.Set("__VIEWSTATEGENERATOR", "3A8BE513")
	logParams.Set("txtDate", logDate)
	logParams.Set("txtStartTime", startTime)
	logParams.Set("txtEndTime", endTime)
	logParams.Set("ddlProjectList", e.projectID)
	logParams.Set("hplbWorkType", "0106")
	logParams.Set("hplbAction", "010601")
	logParams.Set("TextBox1", "")
	//logParams.Set("txtMemo", e.lc.DailyWorkContent)
	logParams.Set("txtMemo", e.workReport.workPlanDaily())
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
	_, err := e.doPost(workLogUrl, strings.NewReader(logParams.Encode()))
	return err
}

// doWorkLogSameDaily 上午、下午的日志内容相同
func (e EDSLogger) doWorkLogSameDaily(workLogUrl, logDate, workPlan string, dt dayTime, hiddenParams map[string]string) error {
	startTime, endTime := dt.startTime, dt.endTime

	logParams := url.Values{}
	logParams.Set("__EVENTTARGET", "hplbWorkType")
	logParams.Set("__EVENTARGUMENT", "")
	logParams.Set("__LASTFOCUS", "")
	logParams.Set("__VIEWSTATEGENERATOR", "3A8BE513")
	logParams.Set("txtDate", logDate)
	logParams.Set("txtStartTime", startTime)
	logParams.Set("txtEndTime", endTime)
	logParams.Set("ddlProjectList", e.projectID)
	logParams.Set("hplbWorkType", "0106")
	logParams.Set("hplbAction", "010601")
	logParams.Set("TextBox1", "")
	logParams.Set("txtMemo", workPlan)
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
	_, err := e.doPost(workLogUrl, strings.NewReader(logParams.Encode()))
	return err
}

func (e EDSLogger) doWeeklyLog(monday string) error {
	logParams := url.Values{}
	logParams.Set("hidCurrRole", "")
	logParams.Set("hidWeeklyState", "")
	logParams.Set("WeekReportDate", monday)
	//logParams.Set("txtWorkContent", e.lc.WeeklyWorkContent)
	logParams.Set("txtWorkContent", e.workReport.LastWeekWorkContent)
	//logParams.Set("txtStudyContent", e.lc.WeeklyStudyContent)
	logParams.Set("txtStudyContent", e.workReport.LastWeekStudyContent)
	//logParams.Set("txtSummary", e.lc.WeeklySummary)
	logParams.Set("txtSummary", e.workReport.LastWeekSummary)
	//logParams.Set("txtPlanWork", e.lc.WeeklyPlanWork)
	logParams.Set("txtPlanWork", e.workReport.workPlanWeekly())
	//logParams.Set("txtPlanStudy", e.lc.WeeklyPlanStudy)
	logParams.Set("txtPlanStudy", e.workReport.StudyPlan)
	logParams.Set("btnSubmit", "%E6%8F%90%E4%BA%A4")

	// 通过get获取一些隐藏参数，用作后台校验
	logUrl := e.urls["weekly"]
	hiddenParams, err := getHiddenParams(logUrl)
	if err != nil {
		return err
	}
	for key, value := range hiddenParams {
		logParams.Set(key, value)
	}

	_, err = e.doPost(logUrl, strings.NewReader(logParams.Encode()))
	if err != nil {
		return err
	}

	log.Println("周报填写成功", monday)
	time.Sleep(2 * time.Second)

	return nil
}

func getHiddenParams(getUrl string) (map[string]string, error) {
	result := make(map[string]string)

	respHtml, err := doGet(getUrl)
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

func WeeklyLog() error {
	//return lg.WeeklyLog()
	return lg.WeeklyLogWithGodateChinese()
}
func (e EDSLogger) WeeklyLogWithGodateChinese() error {
	// 填周报
	// 还是要取当周的工作日，因为不一定都在周一执行，如服务器故障等
	today := godate.Today()

	//fmt.Println("cookie:", e.cookie)
	//return nil

	// 先写周报
	// 只能填写本周周报（周一）!!!
	monday := today.Workdays()[0]
	err := e.doWeeklyLog(monday.String())
	if err != nil {
		return err
	}

	// 填日报
	// 直接填7天日报
	for i := 0; i < 7; i++ {
		date, _ := monday.AddDay(i)
		if chinese.IsWorkDayInChina(date) {
			err = e.DailyLog(date.String())
			if err != nil {
				log.Println("填日报失败:", date, err)
				return err
			} else {
				log.Println("填日报成功:", date)
			}
		} else {
			log.Println(date, "放假")
		}
		time.Sleep(time.Second * 2)
	}

	return nil
}

func (e EDSLogger) WeeklyLog() error {
	today := godate.Today()
	workdays := today.Workdays()

	//fmt.Println("cookie:", e.cookie)
	//return nil

	// 先写周报
	// 只能填写本周周报（周一）!!!
	err := e.doWeeklyLog(workdays[0].String())
	if err != nil {
		return err
	}

	// 再写日报
	for _, day := range workdays {
		err = e.DailyLog(day.String())
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
			err = e.DailyLog(dd)
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
