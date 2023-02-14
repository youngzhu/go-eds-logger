/*
Copyright © 2023 youngzy
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"goeds/logger"
	"log"
)

// weeklyCmd represents the weekly command
var weeklyCmd = &cobra.Command{
	Use:     "weekly(w)",
	Short:   "填写这一周的日志",
	Aliases: []string{"w"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取参数
		urlDaily := viper.GetString("urls.daily")
		urlWeekly := viper.GetString("urls.weekly")
		projectID := viper.GetString("projectID")

		return weeklyAction(urlWeekly, urlDaily, projectID)
	},
}

func weeklyAction(urlWeekly, urlDaily, projectID string) error {
	log.Println("urlWeekly:", urlWeekly)
	log.Println("urlDaily:", urlDaily)
	log.Println("WeeklyWorkContent:", logContent.WeeklyWorkContent)

	return logger.DoWeekly(urlWeekly, urlDaily, projectID, logContent)
}

func init() {
	logCmd.AddCommand(weeklyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// weeklyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// weeklyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
