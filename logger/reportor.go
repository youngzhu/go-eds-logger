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

	var addBody = AddBody{
		Id:            "011a8c96e68471a0cde20c4805f09198",
		DepId:         "35",
		DepartmentId:  "68cfbca7-f4be-11ee-89b1-fa163ea58b38",
		ThirdDepId:    "35",
		WorkFrom:      "08:30",
		WorkTo:        "16:30",
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

	request.Header.Set("Authorization", "eyJhbGciOiJIUzUxMiJ9.eyJjbGllbnQ6bG9naW5fdXNlcl9rZXkiOiI4ODg0ZDNmZS1lMjIzLTQ0NGItYmZhZC0yZGI4YmRmOGM2MzIifQ.pMuI_03uJHf1mdQDYvlajkY0TJawXDyNSDKOQjrH4vIyFj4l15vD8H2ez04Wrj4TyLO0bq-My-W_ZyoUoVYfKQ")
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
