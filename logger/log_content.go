package logger

import (
	"encoding/json"
	"log"
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
