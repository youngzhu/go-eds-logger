package main

import (
	"edser/logger"
	"github.com/youngzhu/godate"
	"github.com/youngzhu/godate/chinese"
	"log"
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	//os.Setenv("EDS_USR_ID", "123")

	log.Println(os.LookupEnv("EDS_USR_ID"))
	log.Println(os.LookupEnv("EDS_USR_PWD"))

	log.Println(os.Getenv("EDS_USR_ID"))
	log.Println(os.Getenv("EDS_USR_PWD"))
}

func TestGetProjectID(t *testing.T) {
	t.Log("projectID:", logger.RetrieveProjectID(cfg))
}

func TestRetrieveLogContent(t *testing.T) {
	days := logger.RetrieveExtraDays()

	t.Log(days)
}

func TestRetrieveHplb(t *testing.T) {
	hplb := logger.RetrieveHplb(cfg)
	t.Log("WorkType:", hplb.WorkType, "Action:", hplb.Action)
}

func TestGodate(t *testing.T) {

	if chinese.IsOffDayInChina(godate.MustDate(2025, 1, 1)) != true {
		t.Fail()
	}

	if chinese.IsOffDayInChina(godate.MustDate(2025, 10, 1)) == true {
		t.Fail()
	}
}
