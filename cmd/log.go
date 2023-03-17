/*
Copyright © 2023 youngzy
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"goeds/logger"
	"path/filepath"
)

var logContent logger.LogContent

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Do EDS log",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// 获取参数
		loggerFilePath := viper.GetString("logger-file")
		if loggerFilePath == "" {
			// Find home directory.
			home, err := homedir.Dir()
			cobra.CheckErr(err)

			loggerFilePath = filepath.Join(home, "edsLogger.json")
		}

		err = logger.RetrieveLogContent(loggerFilePath)
		if err != nil {
			return err
		}

		// 登录
		// 获取参数
		logger.AddUrl("login", viper.GetString("urls.login"))
		logger.SetCookie(viper.GetString("cookie"))
		logger.SetHost(viper.GetString("host"))
		userID := viper.GetString("usr-id")
		userPwd := viper.GetString("usr-pwd")
		err = logger.Login(userID, userPwd)
		if err != nil {
			return err
		}

		// 获取项目编号
		logger.AddUrl("daily", viper.GetString("urls.daily"))
		err = logger.RetrieveProjectID()

		// 给logger添加其他配置
		logger.AddUrl("home", viper.GetString("urls.home"))
		logger.AddUrl("weekly", viper.GetString("urls.weekly"))

		return err
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
