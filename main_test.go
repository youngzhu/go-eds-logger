package main

import (
	"edser/logger"
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
