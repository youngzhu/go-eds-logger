package logger

import (
	"edser/config"
	myhttp "edser/http"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

func RetrieveProjectID(cfg config.Configuration) string {
	if !logined {
		login(cfg)
	}

	respHtml, _ := myhttp.DoGet(cfg.GetStringDefault("urls:worklog", ""))
	//println(respHtml)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(respHtml))

	if err != nil {
		log.Fatalln(err)
	}

	var projectId = ""
	doc.Find("select").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		if id == "ddlProjectList" {
			projectId, _ = s.Children().Attr("value")
			return
		}
	})
	return projectId
}
