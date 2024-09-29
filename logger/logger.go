package logger

import (
	"edser/config"
	"edser/http"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/youngzhu/godate"
	"github.com/youngzhu/godate/chinese"

	//"github.com/youngzhu/godate/chinese"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

var d godate.Date

type EDSLogger interface {
	Execute(cfg config.Configuration) error
}

var loggers = make(map[string]EDSLogger)

func Run(cfg config.Configuration) (err error) {
	// 不测试了，反正也经常失败
	// 检验网站是否正常
	// 	err = http.CheckURL(cfg.GetStringDefault("urls:home", ""))
	// 	if err != nil {
	// 		return err
	// 	}

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
		edsLogger = loggers["Action"]
	} else {
		//return errors.New("手动执行，抛出异常")
	}
	err = edsLogger.Execute(cfg)

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

// 登录
func login(cfg config.Configuration) error {
	user, err := getUser()
	if err != nil {
		return err
	}

	loginURL := cfg.GetStringDefault("urls:login", "")
	err = http.Login(loginURL, user.id, user.password)
	if err != nil {
		return err
	}

	logined = true
	log.Println("登陆成功")

	return nil
}

var (
	ddlProjectList string

	hplbInfo hplb
)

func prepareData(cfg config.Configuration) error {
	ddlProjectList = RetrieveProjectID(cfg)

	hplbInfo = RetrieveHplb(cfg)

	return nil
}

type dayTime struct {
	startTime string
	endTime   string
}

var am = dayTime{startTime: "10:00", endTime: "12:00"}
var pm = dayTime{startTime: "13:00", endTime: "18:00"}

func workLog(cfg config.Configuration, logDate string) (err error) {
	url := cfg.GetStringDefault("urls:worklog", "") + "&LogDate=" + logDate

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams, err := getHiddenParams(url)
	if err != nil {
		return err
	}
	//fmt.Println(len(hiddenParams))

	logContent := cfg.GetStringDefault("logContent:dailyWorkContent", "")
	for _, t := range []dayTime{am, pm} {
		err = doWorkLog(url, logDate, logContent, t, hiddenParams)
		if err != nil {
			return err
		}
	}

	log.Println("日志操作成功", logDate)
	time.Sleep(800 * time.Millisecond)

	return
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
	logParams.Set("hplbWorkType", hplbInfo.WorkType)
	logParams.Set("hplbAction", hplbInfo.Action)
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

	_, err := http.DoPost(workLogUrl, strings.NewReader(logParams.Encode()))
	return err
}

func workWeeklyLog(cfg config.Configuration, logDate string) (err error) {
	urlWeekly := cfg.GetStringDefault("urls:workWeekly", "")

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams, err := getHiddenParams(urlWeekly)
	if err != nil {
		return err
	}

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

	_, err = http.DoPost(urlWeekly, strings.NewReader(logParams.Encode()))
	if err != nil {
		return err
	}

	// resp := http.DoRequest(urlWeekly, http.MethodPost, cookie, strings.NewReader(logParams.Encode()))
	// log.Println(resp)

	log.Println("周报填写成功", logDate)
	time.Sleep(2 * time.Second)

	return
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

func logWholeWeek(cfg config.Configuration, d godate.Date) (err error) {
	workdays := d.Workdays()

	monday := workdays[0]

	// 先写周报
	// 只能填写本周周报（周一）!!!
	err = workWeeklyLog(cfg, monday.String())

	// 再写日报
	// 直接填7天日报
	for i := 0; i < 7; i++ {
		date, _ := d.AddDay(i)
		if chinese.IsWorkDayInChina(date) {
			workLog(cfg, date.String())
		} else {
			log.Println(date, "放假")
		}
	}

	return
}

func register(logType string, edsLogger EDSLogger) {
	if _, exists := loggers[logType]; exists {
		log.Fatalln(logType, "already registered")
	}

	log.Println("Register logger:", logType)
	loggers[logType] = edsLogger
}
