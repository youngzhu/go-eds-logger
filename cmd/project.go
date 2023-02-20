/*
Copyright © 2023 youngzy
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"goeds/logger"
	"log"
	"strings"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:     "project(p)",
	Short:   "获取项目编号",
	Aliases: []string{"p"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取参数
		reqUrl := viper.GetString("urls.daily")

		return projectAction(reqUrl)
	},
}

func projectAction(reqUrl string) error {
	respHtml, _ := logger.DoGet(reqUrl)
	//println(respHtml)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(respHtml))

	if err != nil {
		return err
	}

	var projectId = ""
	doc.Find("select").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		if id == "ddlProjectList" {
			projectId, _ = s.Children().Attr("value")
			return
		}
	})

	log.Println("项目编号：", projectId)

	return nil
}

func init() {
	rootCmd.AddCommand(projectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
