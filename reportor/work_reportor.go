package reportor

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strings"
	"time"
)

// 学 viper 设置一个影子变量
var _reportor *WorkReportor

func init() {
	_reportor = New()
}

type WorkReportor struct {
	workReport WorkReport // 周报内容
	projectId  string
}

func New() *WorkReportor {
	wr := new(WorkReportor)

	//wr.urls = make(map[string]string)

	return wr
}

// 登录
func login() (err error) {
	return _reportor.login()
}

func (r *WorkReportor) login() (err error) {
	userID := viper.GetString("usr-id")
	passcode := viper.GetString("usr-pwd")
	if userID == "" || passcode == "" {
		return errors.New("用户名或密码不能为空")
	}

	params := url.Values{}
	params.Set("UserId", userID)
	params.Set("UserPsd", passcode)

	loginUrl := viper.GetString("urls.login")
	log.Println(loginUrl, "登录中...")
	resp, err := r.doPost(loginUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("登录错误：%w", err)
	}

	//log.Println(resp)

	if strings.Contains(resp, ErrInvalidUser.Error()) {
		return ErrInvalidUser
	}

	//time.Sleep(2 * time.Second)
	log.Println("登陆成功")

	return
}

// 装载周报内容
func loadWorkReport() {
	_reportor.loadWorkReport()
}

func (r *WorkReportor) loadWorkReport() {
	r.workReport = retrieveWorkReport()
}

// 填周报
func fillWeeklyReport(reportDate string) (err error) {
	return _reportor.fillWeeklyReport(reportDate)
}

func (r WorkReportor) fillWeeklyReport(reportDate string) (err error) {
	reportUrl := viper.GetString("urls.weekly")
	// 只能填写本周周报（周一）!!!

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams, err := getHiddenParams(reportUrl)
	if err != nil {
		return err
	}

	// 准备周报的请求参数
	logParams := url.Values{}
	logParams.Set("hidCurrRole", "")
	logParams.Set("hidWeeklyState", "")
	logParams.Set("WeekReportDate", reportDate)
	logParams.Set("txtWorkContent", r.workReport.LastWeekWorkContent)
	logParams.Set("txtStudyContent", r.workReport.StudyPlan)
	logParams.Set("txtSummary", r.workReport.LastWeekSummary)
	logParams.Set("txtPlanWork", r.workReport.workPlanWeekly())
	logParams.Set("txtPlanStudy", r.workReport.StudyPlan)
	logParams.Set("btnSubmit", "%E6%8F%90%E4%BA%A4")

	// 加上校验参数
	for key, value := range hiddenParams {
		logParams.Set(key, value)
	}

	_, err = r.doPost(reportUrl, strings.NewReader(logParams.Encode()))
	if err != nil {
		return err
	}

	log.Println("周报填写成功", reportDate)
	time.Sleep(2 * time.Second)

	return
}

// 填写日报
func fillDailyReport(reportDate string) (err error) {
	return _reportor.fillDailyReport(reportDate)
}

type dayTime struct {
	startTime string
	endTime   string
}

var (
	am = dayTime{startTime: "10:00", endTime: "12:00"}
	pm = dayTime{startTime: "13:00", endTime: "18:00"}
)

func (r WorkReportor) fillDailyReport(reportDate string) (err error) {
	reportUrl := fmt.Sprintf("%s&LogDate=%s", viper.GetString("urls.daily"), reportDate)

	// 先通过get获取一些隐藏参数，用作后台校验
	hiddenParams, err := getHiddenParams(reportUrl)
	if err != nil {
		return fmt.Errorf("获取隐藏参数失败：%w", err)
	}

	for _, t := range []dayTime{am, pm} {
		err := r.fillDailyReportAMPM(reportUrl, reportDate, t, hiddenParams)
		if err != nil {
			return fmt.Errorf("日志操作失败：%w", err)
		}
	}

	//log.Println("日志操作成功", reportDate)
	time.Sleep(time.Second)

	return
}

func (r WorkReportor) fillDailyReportAMPM(reportUrl, reportDate string, dt dayTime, hiddenParams map[string]string) error {
	startTime, endTime := dt.startTime, dt.endTime

	logParams := url.Values{}
	logParams.Set("__EVENTTARGET", "hplbWorkType")
	logParams.Set("__EVENTARGUMENT", "")
	logParams.Set("__LASTFOCUS", "")
	logParams.Set("__VIEWSTATEGENERATOR", "3A8BE513")
	logParams.Set("txtDate", reportDate)
	logParams.Set("txtStartTime", startTime)
	logParams.Set("txtEndTime", endTime)
	logParams.Set("ddlProjectList", r.getProjectId())
	logParams.Set("hplbWorkType", "0106")
	logParams.Set("hplbAction", "010601")
	logParams.Set("TextBox1", "")
	logParams.Set("txtMemo", r.workReport.workPlanDaily())
	logParams.Set("btnSave", "+%E7%A1%AE+%E5%AE%9A+")
	logParams.Set("txtnodate", reportDate)
	logParams.Set("txtnoStartTime", startTime)
	logParams.Set("txtnoEndTime", endTime)
	logParams.Set("TextBox6", "")
	logParams.Set("txtnoMemo", "")
	logParams.Set("txtCRMDate", reportDate)
	logParams.Set("txtCRMStartTime", startTime)
	logParams.Set("txtCRMEndTime", endTime)
	logParams.Set("TextBox5", "")
	logParams.Set("txtCRMMemo", "")

	_, err := r.doPost(reportUrl, strings.NewReader(logParams.Encode()))
	return err
}

func (r *WorkReportor) getProjectId() string {
	if r.projectId == "" {
		// 从页面上获取项目ID
		respHtml, _ := r.doGet(viper.GetString("urls.daily"))

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(respHtml))

		if err != nil {
			log.Fatalln(err)
		}

		var projectId string
		doc.Find("select").Each(func(i int, s *goquery.Selection) {
			id, _ := s.Attr("id")
			if id == "ddlProjectList" {
				projectId, _ = s.Children().Attr("value")
				return
			}
		})

		r.projectId = projectId
	}

	return r.projectId
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
