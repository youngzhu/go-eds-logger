package reportor

import (
	"encoding/json"
	"fmt"
	"github.com/youngzhu/godate"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

//
// 用于提供工作报告的内容
//

type WorkReport struct {
	//上周工作任务完成情况
	LastWeekWorkContent string `json:"lastWeekWorkContent"`
	//上周学习完成任务情况
	LastWeekStudyContent string `json:"lastWeekStudyContent"`
	// 经验和收获总结
	LastWeekSummary string `json:"lastWeekSummary"`
	//本周工作计划与重点
	WorkPlan []string `json:"workPlan"`
	//本周学习计划
	StudyPlan string `json:"studyPlan"`
}

// 以文本形式返回本周工作计划
func (r WorkReport) workPlanWeekly() string {
	return strings.Join(r.WorkPlan, "\n")
}

// 日报
// 从工作计划中随机获取一条
func (r WorkReport) workPlanDaily() string {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机索引
	randomIndex := rand.Intn(len(r.WorkPlan))
	
	return r.WorkPlan[randomIndex]
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
