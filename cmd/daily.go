/*
Copyright © 2023 youngzy
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/youngzhu/godate"
	"goeds/logger"
)

var (
	logStart int
	logDays  int
)

// dailyCmd represents the daily command
// goeds log d -s=0 -d=1
var dailyCmd = &cobra.Command{
	Use:     "d -s=0 -d=1",
	Aliases: []string{"d"},
	Short:   "按天填写日志",
	//Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取参数
		//from, err := strconv.Atoi(args[0])
		//if err != nil {
		//	return err
		//}
		//days, err := strconv.Atoi(args[1])
		//if err != nil {
		//	return err
		//}

		return dailyAction()
	},
}

func dailyAction() error {
	//log.Println(logStart + 1)
	//log.Println(logDays + 1)

	var logDay godate.Date
	for i := 0; i < logDays; i++ {
		diff := logStart + i
		if diff >= 0 {
			logDay, _ = godate.Today().AddDay(diff)
		} else {
			logDay, _ = godate.Today().SubDay(-diff)
		}

		err := logger.DailyLog(logDay.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	logCmd.AddCommand(dailyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dailyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dailyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	dailyCmd.Flags().IntVarP(&logStart, "start", "s", 0, "日志开始的起点。0：今天；-1：昨天；1：明天")
	dailyCmd.Flags().IntVarP(&logDays, "days", "d", 0, "连续填写日志的天数")

}
