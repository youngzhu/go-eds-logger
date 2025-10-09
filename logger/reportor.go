package logger

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

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

func DailyReport(logDate string) error {
	return lg.DailyReport(logDate)
}

func (e EDSLogger) DailyReport(logDate string) error {
	logUrl := e.urls["daily"]
	logUrl = viper.GetString("reportUrl")

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
	addBody.WorkDesc1 = e.workReport.workPlanDaily()

	_, err := doPost(logUrl, addBody)
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

func doPost(url string, entry interface{}) ([]byte, error) {
	entryJson, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}
	return doRequest(url, http.MethodPost, strings.NewReader(string(entryJson)))
}

//type postEntry interface {
//	Body() io.Reader
//}

func doRequest(url, method string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 很重要，代替了以前的安全校验
	request.Header.Set("Authorization", "eyJhbGciOiJIUzUxMiJ9.eyJjbGllbnQ6bG9naW5fdXNlcl9rZXkiOiJmM2Q5ZTJmNy01NTYwLTQzYTEtYjcxNi05MjYzOGFmYjIzYWEifQ.qYQZUJ-yDoo95RtZ2BCET1LBJ-0KmJi1WshQY9aCFcs5UyPejkLukqt0xN-wPm56MynAKSyX-iQuIJXZShtDPA")

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
