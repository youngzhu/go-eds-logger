package logger

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

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

const edspyRoot = "E:\\workspace\\GitHub\\edspy"

// LoadWorkReportRandomly 从本地 edspy 项目中随机加载一个周报文件
func LoadWorkReportRandomly() error {
	return r.LoadWorkReportRandomly()
}

func (re *Reportor) LoadWorkReportRandomly() error {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 指定目录路径
	dirPath := filepath.Join(edspyRoot, "data")

	// 读取目录中的所有文件
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("读取目录失败: %v", err)
	}

	// 过滤出.json文件
	var jsonFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			jsonFiles = append(jsonFiles, file.Name())
		}
	}

	// 检查是否有JSON文件
	if len(jsonFiles) == 0 {
		log.Fatal("目录中没有JSON文件")
	}

	// 随机选择一个JSON文件
	randomFile := jsonFiles[rand.Intn(len(jsonFiles))]
	filePath := filepath.Join(dirPath, randomFile)

	log.Printf("随机选择的文件: %s\n", randomFile)

	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	var workReport WorkReport
	if err := json.Unmarshal(content, &workReport); err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}

	re.workReport = workReport

	//log.Printf("读取到的工作报告: %+v\n", workReport)

	return nil
}
