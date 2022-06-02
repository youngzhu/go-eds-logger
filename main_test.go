package main

import (
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
