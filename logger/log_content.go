package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const defaultPath = "data/logger_default.json"

type LogContent struct {
	DailyWorkContent   string `json:"dailyWorkContent"`
	WeeklyWorkContent  string `json:"weeklyWorkContent"`
	WeeklyStudyContent string `json:"weeklyStudyContent"`
	WeeklySummary      string `json:"weeklySummary"`
	WeeklyPlanWork     string `json:"weeklyPlanWork"`
	WeeklyPlanStudy    string `json:"weeklyPlanStudy"`
}

func RetrieveLogContent(path string) error {
	return lg.RetrieveLogContent(path)
}
func (e *EDSLogger) RetrieveLogContent(path string) error {
	if path != "" {
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("[%s] not exists\n", path)
				path = ""
			}
		}
	}

	if path == "" {
		path = defaultPath
	}

	log.Printf("加载日志内容[%s]...", path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	var content LogContent
	err = json.NewDecoder(file).Decode(&content)
	if err != nil {
		return err
	}

	e.lc = content

	return nil
}

const cdn = "https://fastly.jsdelivr.net"

func RetrieveLogContentViaWeb() error {
	return lg.RetrieveLogContentViaWeb()
}
func (e *EDSLogger) RetrieveLogContentViaWeb() error {
	// 定义要获取的URL
	//url := "https://cdn.jsdelivr.net/gh/youngzhu/edspy/data/work-report-2025-08-11.json"
	//url = "https://github.com/youngzhu/edspy/blob/main/data/work-report-2025-08-11.json"
	url := fmt.Sprintf("%s/gh/youngzhu/edspy/data/work-report-2026-02-09.json", cdn)

	// 发送HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("请求返回非200状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应体失败: %v", err)
	}

	// 解析JSON到空接口（适用于未知结构的JSON）
	//var data interface{}
	//if err := json.Unmarshal(body, &data); err != nil {
	//	log.Fatalf("解析JSON失败: %v", err)
	//}
	var data LogContent
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}

	// 打印获取到的JSON内容
	log.Println("获取到的JSON内容:", data)

	// 如果你想将JSON格式化输出
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatalf("格式化JSON失败: %v", err)
	}
	fmt.Println("\n格式化后的JSON:")
	fmt.Println(string(prettyJSON))

	//e.lc = content

	return nil
}
