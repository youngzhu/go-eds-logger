package logger

import (
	"encoding/json"
	"log"
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

var lc logContent

func loadLogContent() (err error) {
	file, err := os.Open(dataFile)
	if err == nil {
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&lc)
	}

	return
}

func init() {
	err := loadLogContent()
	if err != nil {
		log.Fatal(err)
	}
}
