package logger

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/youngzhu/go-eds-logger/config"
	"github.com/youngzhu/go-eds-logger/http"
	"github.com/youngzhu/godate"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

var d godate.Date

type EDSLogger interface {
	Execute(cfg config.Configuration)
}

var loggers = make(map[string]EDSLogger)

func Run(cfg config.Configuration) (err error) {
	// 登录
	err = login(cfg)
	if err != nil {
		return err
	}

	// 动态获取一些参数
	err = prepareData(cfg)
	if err != nil {
		return err
	}

	//
	edsLogger, exists := loggers["manual"]
	if !exists {
		edsLogger = loggers["action"]
	}
	edsLogger.Execute(cfg)

	return
}

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

var logined = false

func login(cfg config.Configuration) error {
	user, err := getUser()
	if err != nil {
		return err
	}

	// 检验网站是否正常
	err = http.Connect(cfg.GetStringDefault("urls:home", ""))
	if err != nil {
		return err
	}

	// 登录
	// data := `{"UserId":"###", "UserPsd":"***"}`
	// data := "UserId=###&UserPsd=***"
	// params := url.Values{
	// 	"UserId":  {"###"},
	// 	"UserPsd": {"***"},
	// }
	params := url.Values{}
	params.Add("UserId", user.id)
	params.Add("UserPsd", user.password)
	// var request *http.Request
	// request, err = http.NewRequest(http.MethodPost, URL_LOGIN, strings.NewReader(data))
	// request, err = http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(params.Encode()))

	respStr := http.DoPost(cfg.GetStringDefault("urls:login", ""),
		strings.NewReader(params.Encode()))

	// log.Println(respStr)

	errMsg := "用户名或密码错误"
	if strings.Contains(respStr, errMsg) {
		return errors.New(errMsg)
	}

	logined = true
	log.Println("登陆成功")

	return nil
}

var (
	ddlProjectList string
)

func prepareData(cfg config.Configuration) error {
	ddlProjectList = RetrieveProjectID(cfg)

	return nil
}

type dayTime struct {
	startTime string
	endTime   string
}

var am = dayTime{startTime: "10:00", endTime: "12:00"}
var pm = dayTime{startTime: "13:00", endTime: "18:00"}

func workLog(cfg config.Configuration, logDate string) {
	url := cfg.GetStringDefault("urls:worklog", "") + "&LogDate=" + logDate

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams := getHiddenParams(url)
	// fmt.Println(hiddenParams)

	logContent := cfg.GetStringDefault("logContent:dailyWorkContent", "")
	for _, t := range []dayTime{am, pm} {
		doWorkLog(url, logDate, logContent, t, hiddenParams)
	}

	log.Println("日志操作成功", logDate)
	time.Sleep(800 * time.Millisecond)
}

func doWorkLog(workLogUrl, logDate, logContent string, dt dayTime, hiddenParams map[string]string) {
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

	http.DoPost(workLogUrl, strings.NewReader(logParams.Encode()))

}

func workWeeklyLog(cfg config.Configuration, logDate string) {
	urlWeekly := cfg.GetStringDefault("urls:workWeekly", "")

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams := getHiddenParams(urlWeekly)

	logParams := url.Values{}
	logParams.Set("hidCurrRole", "")
	logParams.Set("hidWeeklyState", "")
	logParams.Set("WeekReportDate", logDate)
	logParams.Set("txtWorkContent",
		cfg.GetStringDefault("logContent:weeklyWorkContent", ""))
	logParams.Set("txtStudyContent",
		cfg.GetStringDefault("logContent:weeklyStudyContent", ""))
	logParams.Set("txtSummary",
		cfg.GetStringDefault("logContent:weeklySummary", ""))
	logParams.Set("txtPlanWork",
		cfg.GetStringDefault("logContent:weeklyPlanWork", ""))
	logParams.Set("txtPlanStudy",
		cfg.GetStringDefault("logContent:weeklyPlanStudy", ""))
	logParams.Set("btnSubmit", "%E6%8F%90%E4%BA%A4")

	for key, value := range hiddenParams {
		logParams.Set(key, value)
	}

	http.DoPost(urlWeekly, strings.NewReader(logParams.Encode()))

	// resp := http.DoRequest(urlWeekly, http.MethodPost, cookie, strings.NewReader(logParams.Encode()))
	// log.Println(resp)

	log.Println("周报填写成功", logDate)
	time.Sleep(2 * time.Second)
}

func getHiddenParams(url string) map[string]string {
	result := make(map[string]string)

	respHtml := http.DoGet(url)
	//println(respHtml)

	keys := []string{"__EVENTVALIDATION", "__VIEWSTATE"}

	for _, k := range keys {
		result[k] = getValueFromHtml(respHtml, k)
	}

	return result
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

func logWholeWeek(cfg config.Configuration, d godate.Date) {
	workdays := d.Workdays()

	// 先写周报
	// 只能填写本周周报（周一）!!!
	workWeeklyLog(cfg, workdays[0].String())

	// 再写日报
	for _, day := range workdays {
		workLog(cfg, day.String())
	}

	// 周末调休
	sat, _ := d.AddDay(5)
	sun, _ := d.AddDay(6)

	extraDays := RetrieveExtraDays()

	for _, dd := range []string{sat.String(), sun.String()} {
		if _, ok := extraDays[dd]; ok {
			log.Println("调休", dd)
			workLog(cfg, dd)
		}
	}
}

func register(logType string, edsLogger EDSLogger) {
	if _, exists := loggers[logType]; exists {
		log.Fatalln(logType, "already registered")
	}

	log.Println("Register logger:", logType)
	loggers[logType] = edsLogger
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
