package logger

import (
	"encoding/json"
	"errors"
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

func LoginX(userId, password string) error {
	return r.Login(userId, password)
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
	var loginBody = LoginReq{
		LoginType:      "password",
		EnterpriseCode: "Newtouch",
		EmployeeId:     userId,
		Password:       password,
	}
	loginUrl := viper.GetString("loginUrl")
	loginUrl = "https://eds.newtouch.com/api/login"

	resp, err := re.postJSON(loginUrl, loginBody)
	if err != nil {
		return err
	}

	var loginResp LoginResp
	err = json.Unmarshal(resp, &loginResp)
	if err != nil {
		return err
	}
	if loginResp.Code != 200 {
		return errors.New("登录失败: " + loginResp.Msg)
	}

	re.Token = loginResp.Token

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
func DailyReport(logDate string) error {
	return r.DailyReport(logDate)
}

func (re Reportor) DailyReport(logDate string) error {
	logUrl := viper.GetString("reportUrl")

	logUrl = "https://eds.newtouch.com/api/workReport/add"
	if logUrl == "" {
		return errors.New("logUrl为空")
	}

	/*
		{
		    "id": "",
		    "depId": "35",
		    "departmentId": "68cfbca7-f4be-11ee-89b1-fa163ea58b38",
		    "thirdDepId": "35",
		    "reportDate": "2025-10-09",
		    "workFrom": "08:30",
		    "workTo": "16:30",
		    "action1Id": "",
		    "action2Id": "",
		    "actionFirstId": "",
		    "actionSecondId": "",
		    "workDesc1": "投连产品",
		    "workDesc2": "PC",
		    "workHours": 8,
		    "isHaveProject": "有",
		    "proRecordId": 16205,
		    "proId": "Q2503017",
		    "timeType": 0
		}
	*/
	/*
		{
		    "id": "",
		    "depId": "35",
		    "departmentId": "68cfbca7-f4be-11ee-89b1-fa163ea58b38",
		    "thirdDepId": "35",
		    "reportDate": "2025-10-11",
		    "workFrom": "08:30",
		    "workTo": "16:30",
		    "action1Id": "",
		    "action2Id": "",
		    "actionFirstId": "",
		    "actionSecondId": "",
		    "workDesc1": "撤退减保，投连账户优化",
		    "workDesc2": "PC",
		    "workHours": 8,
		    "isHaveProject": "有",
		    "proRecordId": 16205,
		    "proId": "Q2503017",
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
		ProRecordId:   16205,
		ProId:         "Q2503017",
		TimeType:      0,
	}

	addBody.ReportDate = logDate
	addBody.WorkDesc1 = re.workReport.workPlanDaily()

	_, err := re.postJSON(logUrl, addBody)
	if err != nil {
		return err
	}

	log.Println("日志操作成功", logDate)
	time.Sleep(800 * time.Millisecond)

	return nil
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

func WeeklyReport() error {
	return r.WeeklyReport()
}
func (re Reportor) WeeklyReport() error {
	// 填周报
	// 还是要取当周的工作日，因为不一定都在周一执行，如服务器故障等
	today := godate.Today()

	//fmt.Println("cookie:", e.cookie)
	//return nil

	// 先写周报
	// 只能填写本周周报（周一）!!!
	monday := today.Workdays()[0]
	//err := e.doWeeklyLog(monday.String())
	//if err != nil {
	//	return err
	//}

	var err error

	// 填日报
	// 直接填7天日报
	for i := 0; i < 7; i++ {
		date, _ := monday.AddDay(i)
		if chinese.IsWorkDayInChina(date) {
			err = re.DailyReport(date.String())
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
