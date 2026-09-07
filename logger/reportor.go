package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngzhu/godate"
	"github.com/youngzhu/godate/chinese"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// InitConfig 初始化配置
func InitConfig() error {
	// 读取环境变量
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("EDS")
	viper.AutomaticEnv() // read in environment variables that match

	// 读取配置文件
	// Find home directory.
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".goeds" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".goeds")

	// If a config file is found, read it in.
	return viper.ReadInConfig()

}

type Reportor struct {
	Token      string
	workReport WorkReport
}

// 学 viper 设置一个影子变量
var r *Reportor

func init() {
	r = new(Reportor)
}

// 登录的请求和响应体
type (
	LoginReq struct {
		LoginType      string `json:"loginType"`
		EnterpriseCode string `json:"enterpriseCode"`
		EmployeeId     string `json:"employeeId"`
		Password       string `json:"password"`
	}

	LoginResp struct {
		Msg   string `json:"msg"`
		Code  int    `json:"code"`
		Token string `json:"token"`
	}
)

func LoginX() error {
	userID := viper.GetString("usr-id")
	userPwd := viper.GetString("usr-pwd")
	return r.Login(userID, userPwd)
}

func (re *Reportor) Login(userId, password string) error {
	/*
		{
		    "loginType": "password",
		    "enterpriseCode": "Newtouch",
		    "employeeId": "",
		    "password": ""
		}
	*/
	/*
	var loginBody = LoginReq{
		LoginType:      "password",
		EnterpriseCode: "Newtouch",
		EmployeeId:     userId,
		Password:       password,
	}
	loginUrl := viper.GetString("urls.login")

	resp, err := re.postJSON(loginUrl, loginBody)
	if err != nil {
		return err
	}

	var loginResp LoginResp
	err = json.Unmarshal(resp, &loginResp)
	if err != nil {
		return err
	}

	log.Println("登录返回: ", string(resp))
	if loginResp.Code != 200 {
		return errors.New("登录失败: " + loginResp.Msg)
	}
*/

	// 登录返回取消了token
	//re.Token = loginResp.Token
	re.Token = viper.GetString("token")

	log.Println("登陆成功")

	return nil
}

type AddBody struct {
	Id             string `json:"id"`
	DepId          string `json:"depId"`
	DepartmentId   string `json:"departmentId"`
	ThirdDepId     string `json:"thirdDepId"`
	ReportDate     string `json:"reportDate"`
	WorkFrom       string `json:"workFrom"`
	WorkTo         string `json:"workTo"`
	Action1Id      string `json:"action1Id"`
	Action2Id      string `json:"action2Id"`
	ActionFirstId  string `json:"actionFirstId"`
	ActionSecondId string `json:"actionSecondId"`
	WorkDesc1      string `json:"workDesc1"`
	WorkDesc2      string `json:"workDesc2"`
	WorkHours      int    `json:"workHours"`
	IsHaveProject  string `json:"isHaveProject"`
	ProRecordId    int    `json:"proRecordId"`
	ProId          string `json:"proId"`
	TimeType       int    `json:"timeType"`
}

// DailyReport 填日报
// 注意：只填空白的。填过的，不会更新了
func DailyReport(reportDate string) error {
	return r.DailyReport(reportDate)
}

func (re Reportor) DailyReport(reportDate string) error {
	logUrl := viper.GetString("urls.daily")

	if logUrl == "" {
		return errors.New("logUrl为空")
	}

	/*
		{
		    "id": "fd0d433d512b40e297f9938831850033",
		    "depId": "35",
		    "departmentId": "68cfbca7-f4be-11ee-89b1-fa163ea58b38",
		    "thirdDepId": "35",
		    "reportDate": "2026-05-15",
		    "workFrom": "08:30",
		    "workTo": "16:30",
		    "action1Id": "",
		    "action2Id": "",
		    "actionFirstId": "",
		    "actionSecondId": "",
		    "workDesc1": "支持新产品的上线需求，完成业务场景的代码开发。",
		    "workDesc2": "PC",
		    "workHours": 8,
		    "isHaveProject": "有",
		    "proRecordId": 18368,
		    "proId": "18368",
		    "timeType": 0
		}
	*/
	var addBody = AddBody{
		// 不需要，大概是查询用的
		//Id:            "",
		DepId:         "35",
		DepartmentId:  "68cfbca7-f4be-11ee-89b1-fa163ea58b38",
		ThirdDepId:    "35",
		WorkFrom:      "08:30",
		WorkTo:        "16:30",
		WorkDesc2:     "PC",
		WorkHours:     8,
		IsHaveProject: "有",
		//ProRecordId:   18026,
		//ProId:         "Q2602062",
		TimeType: 0,
	}

	addBody.ProRecordId = viper.GetInt("project.recordId")
	addBody.ProId = viper.GetString("project.id")

	addBody.ReportDate = reportDate
	addBody.WorkDesc1 = re.workReport.workPlanDaily()

	_, err := re.postJSON(logUrl, addBody)
	if err != nil {
		return err
	}

	//log.Println("日志操作成功", reportDate)
	//time.Sleep(800 * time.Millisecond)

	return nil
}

type (
	QueryResp struct {
		Msg  string              `json:"msg"`
		Code int                 `json:"code"`
		Data []DailyReportDetail `json:"data"`
	}

	DailyReportDetail struct {
		Id                    string      `json:"id"`
		CstCreate             string      `json:"cstCreate"`
		CreateUserId          string      `json:"createUserId"`
		CstModified           string      `json:"cstModified"`
		UpdateUserId          string      `json:"updateUserId"`
		DeleteUserId          interface{} `json:"deleteUserId"`
		CstDeleted            interface{} `json:"cstDeleted"`
		DeleteFlag            bool        `json:"deleteFlag"`
		Version               int         `json:"version"`
		NewEdsId              int         `json:"newEdsId"`
		EmpId                 string      `json:"empId"`
		EmployeeName          string      `json:"employeeName"`
		EmployeeId            string      `json:"employeeId"`
		DepartmentId          string      `json:"departmentId"`
		CompanyId             string      `json:"companyId"`
		DepId                 string      `json:"depId"`
		ThirdDepId            string      `json:"thirdDepId"`
		ReportDate            string      `json:"reportDate"`
		WorkTimeFrom          string      `json:"workTimeFrom"`
		WorkTimeTo            string      `json:"workTimeTo"`
		Action1Id             string      `json:"action1Id"`
		Action2Id             string      `json:"action2Id"`
		ActionFirstId         string      `json:"actionFirstId"`
		ActionSecondId        string      `json:"actionSecondId"`
		WorkDesc1             string      `json:"workDesc1"`
		WorkDesc2             string      `json:"workDesc2"`
		Sci                   string      `json:"sci"`
		WorkHours             float64     `json:"workHours"`
		AssnUid               interface{} `json:"assnUid"`
		IsHaveProject         string      `json:"isHaveProject"`
		ProId                 string      `json:"proId"`
		ProRecordId           int         `json:"proRecordId"`
		ProjId                interface{} `json:"projId"`
		SbsId                 int         `json:"sbsId"`
		CreateDate            string      `json:"createDate"`
		WorkHours2            interface{} `json:"workHours2"`
		Complainant           interface{} `json:"complainant"`
		TimeType              int         `json:"timeType"`
		ProjName              string      `json:"projName"`
		DepartmentName        interface{} `json:"departmentName"`
		Action1Name           interface{} `json:"action1Name"`
		Action2Name           interface{} `json:"action2Name"`
		IsSettlement          interface{} `json:"isSettlement"`
		ProjectDepartmentName interface{} `json:"projectDepartmentName"`
		LogType               interface{} `json:"logType"`
		LastFillTime          interface{} `json:"lastFillTime"`
		ViewIfOperate         interface{} `json:"viewIfOperate"`
	}
)

func QueryDailyReport(reportDate string) (DailyReportDetail, error) {
	return r.QueryDailyReport(reportDate)
}

func (re Reportor) QueryDailyReport(reportDate string) (DailyReportDetail, error) {
	var detail DailyReportDetail

	queryUrl := viper.GetString("urls.query")
	if queryUrl == "" {
		return detail, errors.New("queryUrl为空")
	}
	queryUrl = fmt.Sprintf(queryUrl, reportDate)

	resp, err := re.doRequest(queryUrl, http.MethodGet, nil)
	if err != nil {
		return detail, err
	}
	//log.Println("查询返回", resp)

	var queryResp QueryResp
	json.NewDecoder(bytes.NewReader(resp)).Decode(&queryResp)
	if queryResp.Code != 200 {
		return detail, errors.New("查询日志失败：" + queryResp.Msg)
	} else {
		if len(queryResp.Data) == 0 {
			return detail, errors.New("未查询到对应日期的日志")
		}
		detail = queryResp.Data[0]
	}

	return detail, nil
}

// HasReport 判断是否有日志
func HasReport(reportDate string) bool {
	return r.HasReport(reportDate)
}

func (re Reportor) HasReport(reportDate string) bool {
	detail, err := re.QueryDailyReport(reportDate)
	if err == nil && detail.Id != "" {
		return true
	}

	return false
}

// DailyReportSafe 保证日志填写成功
// 填完日志查询一下，没有再填一次，最多尝试3次
func DailyReportSafe(reportDate string) {
	// 尝试3次
	for i := 0; i < 3; i++ {
		err := DailyReport(reportDate)
		if err != nil {
			log.Println("日志填写错误", err)
		}

		// 等待0.8秒，避免请求过快被拒绝
		time.Sleep(time.Millisecond * 800)

		// 查询是否填写成功
		if HasReport(reportDate) {
			log.Println("日志填写成功", reportDate)
			return
		} else {
			count := i + 1
			time.Sleep(time.Second * time.Duration(count)) // 等待更久一些
			log.Println("第", count, "次填写失败")
		}
	}

}

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
}

func (re Reportor) postJSON(url string, entry interface{}) ([]byte, error) {
	entryJson, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}
	return re.doRequest(url, http.MethodPost, strings.NewReader(string(entryJson)))
}

//type postEntry interface {
//	Body() io.Reader
//}

func (re Reportor) doRequest(url, method string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 很重要，代替了以前的安全校验
	if re.Token != "" {
		request.Header.Set("Authorization", re.Token)
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//if resp.StatusCode != http.StatusOK {
	//	msg, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		return nil, fmt.Errorf("cannot read body: %w", err)
	//	}
	//	return nil, fmt.Errorf("%w: %s, %s",
	//		err, http.StatusText(resp.StatusCode), msg)
	//}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// WeeklyReport 填写周报和整周的日报
func WeeklyReport() error {
	return r.WeeklyReport()
}
func (re Reportor) WeeklyReport() error {
	// 填周报
	// 还是要取当周的工作日，因为不一定都在周一执行，如服务器故障等
	today := godate.Today()

	// 先写周报
	// 只能填写本周周报（周一）!!!
	monday := chinese.NewCNDate(today.Workdays()[0])
	// 周一是工作日才填周报
	if !monday.Offday {
		err := WeekReport(monday.String())
		if err != nil {
			return err
		}
	}

	// 循环一周填写日报
	date := monday
	for {
		if date.Offday {
			log.Println(date, "放假-", date.Name)
		} else {
			DailyReportSafe(date.String())
		}

		date = date.NextDay()

		// 下周一了，结束
		if date.Weekday() == godate.Monday {
			break
		}
	}

	return nil
}

type WeekReportReq struct {
	WeekReportDate string `json:"weekreportdate"`
	Id             string `json:"id"`
	UnfinishWork   string `json:"unfinishwork"`
	WorkProblem    string `json:"workproblem"`
	Remark         string `json:"remark"`
	Arrangement    string `json:"arrangement"`
	PlanWork       string `json:"planwork"`
	WeekState      string `json:"weekstate"`
}

// WeekReport 填周报
func WeekReport(reportDate string) error {
	return r.WeekReport(reportDate)
}

func (re *Reportor) WeekReport(reportDate string) error {
	/*
		{
		    "weekreportdate": "2025-10-13",
		    "id": "",
		    "unfinishwork": "1 完成一个优化任务\n2 投连需求",
		    "workproblem": "退保挽留通知流程",
		    "remark": "通知流程更熟悉了",
		    "arrangement": "1 代办流程\n2 选卡优化",
		    "planwork": "CodeBuddy",
		    "weekstate": "1"
		}
	*/
	var req = WeekReportReq{
		WeekReportDate: reportDate,
		UnfinishWork:   re.workReport.LastWeekWorkContent,
		WorkProblem:    re.workReport.LastWeekStudyContent,
		Remark:         re.workReport.LastWeekSummary,
		Arrangement:    re.workReport.workPlanWeekly(),
		PlanWork:       re.workReport.StudyPlan,
		WeekState:      "1",
	}

	logUrl := viper.GetString("urls.weekly")

	_, err := re.postJSON(logUrl, req)
	if err != nil {
		return err
	}

	log.Println("周报填写成功", reportDate)
	time.Sleep(2 * time.Second)

	return nil
}
