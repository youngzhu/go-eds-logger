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
	myhttp "goeds/http"
	"log"
	"strings"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "获取项目编号",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取参数
		loginURL := viper.GetString("urls.login")
		userID := viper.GetString("usr-id")
		userPwd := viper.GetString("usr-pwd")

		return projectAction(loginURL, userID, userPwd)
	},
}

func projectAction(loginURL, userID, userPwd string) error {
	err := myhttp.Login(loginURL, userID, userPwd)
	if err != nil {
		return err
	}

	respHtml, _ := myhttp.DoGet(viper.GetString("urls.daily"))
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
