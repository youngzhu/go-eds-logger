package logger

import (
	"github.com/PuerkitoBio/goquery"
	myhttp "github.com/youngzhu/go-eds-logger/http"
	"log"
	"net/http"
	"strings"
)

func RetrieveProjectID() string {
	respHtml := myhttp.DoRequest(workLogURL, http.MethodGet, cookie, nil)
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
