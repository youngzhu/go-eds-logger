package logger

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	myhttp "github.com/youngzhu/go-eds-logger/http"
	"github.com/youngzhu/go-eds-logger/secret"

	"github.com/PuerkitoBio/goquery"
)

const cookie = "ASP.NET_SessionId=4khtnz55xiyhbmncrzmzyzzc; ActionSelect=010601; Hm_lvt_416c770ac83a9d996d7b3793f8c4994d=1569767826; Hm_lpvt_416c770ac83a9d996d7b3793f8c4994d=1569767826; PersonId=12234"

var secretInfo *secret.Secret

// Login 登录
// secretStr 输入参数
// 如果为空，则从文件中读取
func Login(loginInfo *secret.Secret) error {
	// 检验网站是否正常
	resp, err := http.Head(myhttp.UrlHome) // 只请求网站的 http header信息
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if loginInfo != nil {
		secretInfo = loginInfo
	} else {
		secretInfo, err = secret.RetrieveSecret()
		if err != nil {
			log.Fatal(err)
		}
	}

	loginUrl := "http://eds.newtouch.cn/eds3/DefaultLogin.aspx?lan=zh-cn"
	// 登录
	// data := `{"UserId":"###", "UserPsd":"***"}`
	// data := "UserId=###&UserPsd=***"
	// params := url.Values{
	// 	"UserId":  {"###"},
	// 	"UserPsd": {"***"},
	// }
	params := url.Values{}
	params.Add("UserId", secretInfo.UserId)
	params.Add("UserPsd", secretInfo.UserPsd)
	// var request *http.Request
	// request, err = http.NewRequest(http.MethodPost, URL_LOGIN, strings.NewReader(data))
	// request, err = http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(params.Encode()))

	respStr := myhttp.DoRequest(loginUrl, http.MethodPost, cookie, strings.NewReader(params.Encode()))

	// log.Println(respStr)

	errMsg := "用户名或密码错误"
	if strings.Contains(respStr, errMsg) {
		log.Fatalln(errMsg)
	}

	log.Println("登陆成功")

	return nil
}

var workLogURL = "http://eds.newtouch.cn/eds3/worklog.aspx?tabid=0"

func workLog(logDate string) {
	url := workLogURL + "&LogDate=" + logDate

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams := getHiddenParams(url)
	// fmt.Println(hiddenParams)

	doWorkLog(url, logDate, "AM", hiddenParams)
	doWorkLog(url, logDate, "PM", hiddenParams)
}

var dailyLog = "需求的开发、联调与测试" // 日志工作内容

func doWorkLog(workLogUrl, logDate, timeFlag string, hiddenParams map[string]string) {
	startTime := "10:00"
	endTime := "12:00"
	if "PM" == timeFlag {
		startTime = "13:00"
		endTime = "18:00"
	}

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
	logParams.Set("txtMemo", dailyLog)
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

	log.Println("日志操作成功", logDate, timeFlag)
}

var weeklyLog = struct {
	workContent  string
	studyContent string
	summary      string
	planWork     string
	planStudy    string
}{
	workContent:  "需求的讨论/开发/测试与联调",
	studyContent: "JSON的封装与解析",
	summary:      "好的设计，刚开发时可能觉得麻烦似乎没有必要，但再次改动或二次开发时会容易的多",
	planWork:     "新需求的讨论/设计与开发",
	planStudy:    "代码的重构",
}

func workWeeklyLog(logDate string) {
	logUrl := "http://eds.newtouch.cn/eds36web/WorkWeekly/WorkWeeklyInfo.aspx"

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams := getHiddenParams(logUrl)

	logParams := url.Values{}
	logParams.Set("hidCurrRole", "")
	logParams.Set("hidWeeklyState", "")
	logParams.Set("WeekReportDate", logDate)
	logParams.Set("txtWorkContent", weeklyLog.workContent)
	logParams.Set("txtStudyContent", weeklyLog.studyContent)
	logParams.Set("txtSummary", weeklyLog.summary)
	logParams.Set("txtPlanWork", weeklyLog.planWork)
	logParams.Set("txtPlanStudy", weeklyLog.planStudy)
	logParams.Set("btnSubmit", "%E6%8F%90%E4%BA%A4")

	for key, value := range hiddenParams {
		logParams.Set(key, value)
	}

	myhttp.DoRequest(logUrl, http.MethodPost, cookie, strings.NewReader(logParams.Encode()))

	// resp := myhttp.DoRequest(logUrl, http.MethodPost, cookie, strings.NewReader(logParams.Encode()))
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

// LogTheSpecifiedDay 按指定的日期填写日报（只填当天）
func LogTheSpecifiedDay(logDate time.Time) {
	workLog(logDate.Format("2006-01-02"))
}

// LogFromSpecifiedDay 从周一开始，填写本周的周报和日报
func LogFromSpecifiedDay(logFrom time.Time) {

	logDateWeekly := logFrom.Format("2006-01-02")
	logDateDaliy := logFrom

	// 先写周报
	// 只能填写本周周报（周一）!!!
	workWeeklyLog(logDateWeekly)

	time.Sleep(10 * time.Second)

	// 再写日报
	for i := 0; i < 5; i++ {
		LogTheSpecifiedDay(logDateDaliy)

		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
		logDateDaliy = logDateDaliy.Add(time.Hour * 24)
	}

}
