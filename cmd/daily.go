/*
Copyright © 2023 youngzy
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

// dailyCmd represents the daily command
var dailyCmd = &cobra.Command{
	Use:     "daily(d) m n",
	Aliases: []string{"d"},
	Short:   "A brief description of your command",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取参数
		logUrl := viper.GetString("urls.daily")

		from, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		days, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		return dailyAction(logUrl, from, days)
	},
}

func dailyAction(logUrl string, from, days int) error {
	//log.Println(from + 1)
	//log.Println(days + 1)

	for i := 0; i < days; i++ {

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
}
