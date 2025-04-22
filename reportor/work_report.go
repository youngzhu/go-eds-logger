package reportor

import (
	"encoding/json"
	"fmt"
	"github.com/youngzhu/godate"
	"io"
	"log"
	"os"
)

//
// 用于提供工作报告的内容
//

type WorkReport struct {
	LastWeekWorkContent  string   `json:"lastWeekWorkContent"`
	LastWeekStudyContent string   `json:"lastWeekStudyContent"`
	LastWeekSummary      string   `json:"lastWeekSummary"`
	WorkPlan             []string `json:"workPlan"`
	StudyPlan            string   `json:"studyPlan"`
}

// 获取周报和日报的内容
// 1. 首先从 GitHub 上获取
// 2. 如果失败，则从本地获取
func retrieveWorkReport() (workReport WorkReport) {
	log.Println("获取工作报告中...")
	workReport, err := retrieveWorkReportFromInternet()
	if err != nil {
		log.Println("从网络获取工作报告失败，尝试从本地获取：", err.Error())
		workReport = retrieveWorkReportFromLocal()
	}
	log.Println("获取工作报告成功！")
	return
}

func retrieveWorkReportFromLocal() (workReport WorkReport) {
	// Open the file.
	file, err := os.Open("work-report-default.json")
	if err != nil {
		panic(err)
	}

	// Schedule the file to be closed once
	// the function returns.
	defer file.Close()

	workReport, _ = readFromJson(file)

	return
}

func retrieveWorkReportFromInternet() (WorkReport, error) {
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/youngzhu/edspy/data/work-report-%s.json", godate.Today().Workdays()[0])
	url = "https://cdn.jsdelivr.net/gh/youngzhu/edspy/data/work-report-2025-04-16.json"
	get, err := newClient().Get(url)
	//get, err := http.Get(url)
	if err != nil {
		return WorkReport{}, err
	}
	defer get.Body.Close()

	return readFromJson(get.Body)
}

func readFromJson(jsonContent io.Reader) (WorkReport, error) {
	var workReport WorkReport
	err := json.NewDecoder(jsonContent).Decode(&workReport)
	if err != nil {
		return WorkReport{}, err
	}

	return workReport, err
}
