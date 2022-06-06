package logger

import (
	"encoding/json"
	"os"
)

const dataFile = "data/data.json"

type logContent struct {
	DailyWorkContent   string `json:"dailyWorkContent"`
	WeeklyWorkContent  string `json:"weeklyWorkContent"`
	WeeklyStudyContent string `json:"weeklyStudyContent"`
	WeeklySummary      string `json:"weeklySummary"`
	WeeklyPlanWork     string `json:"weeklyPlanWork"`
	WeeklyPlanStudy    string `json:"weeklyPlanStudy"`
}

func RetrieveLogContent() (*logContent, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var content *logContent
	err = json.NewDecoder(file).Decode(&content)

	return content, err
}
