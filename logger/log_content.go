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

func RetrieveLogContent(path string) (LogContent, error) {
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
		log.Println("use", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return LogContent{}, err
	}

	defer file.Close()

	var content LogContent
	err = json.NewDecoder(file).Decode(&content)

	return content, err
}
