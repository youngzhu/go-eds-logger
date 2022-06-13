package logger

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	myhttp "github.com/youngzhu/go-eds-logger/http"
	"github.com/youngzhu/godate"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var d godate.Date

type EDSLogger interface {
	Execute()
}

var loggers = make(map[string]EDSLogger)

func Run() {
	edsLogger, exists := loggers["manual"]
	if !exists {
		edsLogger = loggers["action"]
	}
	edsLogger.Execute()
}

const cookie = "ASP.NET_SessionId=4khtnz55xiyhbmncrzmzyzzc; ActionSelect=010601; Hm_lvt_416c770ac83a9d996d7b3793f8c4994d=1569767826; Hm_lpvt_416c770ac83a9d996d7b3793f8c4994d=1569767826; PersonId=12234"

const loginEnvErr = "环境变量[%s]未配置\n"

type user struct {
	id       string
	password string
}

func getUser() (*user, error) {
	var loginID, loginPassword string
	var found bool
	if loginID, found = os.LookupEnv("EDS_USR_ID"); !found {
		return nil, fmt.Errorf(loginEnvErr, "EDS_USR_ID")
	}

	if loginPassword, found = os.LookupEnv("EDS_USR_PWD"); !found {
		return nil, fmt.Errorf(loginEnvErr, "EDS_USR_PWD")
	}

	return &user{id: loginID, password: loginPassword}, nil
}

const loginUrl = "http://eds.newtouch.cn/eds3/DefaultLogin.aspx?lan=zh-cn"

func Login() error {
	user, err := getUser()
	if err != nil {
		return err
	}

	// 检验网站是否正常
	_, err = http.Head(myhttp.UrlHome) // 只请求网站的 http header信息
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

	respStr := myhttp.DoRequest(loginUrl, http.MethodPost, cookie, strings.NewReader(params.Encode()))

	// log.Println(respStr)

	errMsg := "用户名或密码错误"
	if strings.Contains(respStr, errMsg) {
		return errors.New(errMsg)
	}

	log.Println("登陆成功")

	return nil
}

var (
	lc             *logContent
	ddlProjectList string
)

func PrepareData() error {
	c, err := RetrieveLogContent()
	if err != nil {
		return err
	}
	lc = c

	ddlProjectList = RetrieveProjectID()

	return nil
}

const workLogURL = "http://eds.newtouch.cn/eds3/worklog.aspx?tabid=0"

type dayTime struct {
	startTime string
	endTime   string
}

var am = dayTime{startTime: "10:00", endTime: "12:00"}
var pm = dayTime{startTime: "13:00", endTime: "18:00"}

func workLog(logDate string) {
	url := workLogURL + "&LogDate=" + logDate

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams := getHiddenParams(url)
	// fmt.Println(hiddenParams)

	doWorkLog(url, logDate, am, hiddenParams)
	doWorkLog(url, logDate, pm, hiddenParams)

	log.Println("日志操作成功", logDate)
}

func doWorkLog(workLogUrl, logDate string, dt dayTime, hiddenParams map[string]string) {
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
	logParams.Set("txtMemo", lc.DailyWorkContent)
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

	myhttp.DoRequest(workLogUrl, http.MethodPost, cookie, strings.NewReader(logParams.Encode()))

}

const urlWeekly = "http://eds.newtouch.cn/eds36web/WorkWeekly/WorkWeeklyInfo.aspx"

func workWeeklyLog(logDate string) {

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams := getHiddenParams(urlWeekly)

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

	for key, value := range hiddenParams {
		logParams.Set(key, value)
	}

	myhttp.DoRequest(urlWeekly, http.MethodPost, cookie, strings.NewReader(logParams.Encode()))

	// resp := myhttp.DoRequest(urlWeekly, http.MethodPost, cookie, strings.NewReader(logParams.Encode()))
	// log.Println(resp)

	log.Println("周报填写成功", logDate)
}

func getHiddenParams(url string) map[string]string {
	result := make(map[string]string)

	respHtml := myhttp.DoRequest(url, http.MethodGet, cookie, nil)
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

// 按指定的日期填写日报（只填当天）
func logTheSpecifiedDay(logDate time.Time) {
	//workLog(ParseToStr(logDate))
}

func logWholeWeek(d godate.Date) {
	workdays := d.Workdays()

	// 先写周报
	// 只能填写本周周报（周一）!!!
	workWeeklyLog(workdays[0].String())

	time.Sleep(5 * time.Second)

	// 再写日报
	for _, day := range workdays {
		workLog(day.String())
		time.Sleep(1 * time.Second)
	}
}

func register(logType string, edsLogger EDSLogger) {
	if _, exists := loggers[logType]; exists {
		log.Fatalln(logType, "already registered")
	}

	log.Println("Register logger:", logType)
	loggers[logType] = edsLogger
}
