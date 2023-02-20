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
	"log"
	"path/filepath"
)

var logContent logger.LogContent

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Do EDS log",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 获取参数
		loggerFilePath := viper.GetString("logger-file")
		if loggerFilePath == "" {
			// Find home directory.
			home, err := homedir.Dir()
			cobra.CheckErr(err)

			loggerFilePath = filepath.Join(home, "edsLogger.json")
		}

		err := loadLoggerFile(loggerFilePath)
		if err != nil {
			return err
		}

		// 登录
		// 获取参数
		loginURL := viper.GetString("urls.login")
		userID := viper.GetString("usr-id")
		userPwd := viper.GetString("usr-pwd")

		return logger.Login(loginURL, userID, userPwd)
	},
}

func loadLoggerFile(path string) (err error) {
	log.Printf("加载日志内容[%s]...", path)
	logContent, err = logger.RetrieveLogContent(path)

	return
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
