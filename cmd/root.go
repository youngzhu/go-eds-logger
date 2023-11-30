/*
Copyright © 2023 youngzy
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "goeds",
	Short:         "Do EDS log via Golang",
	Version:       "1.0",
	SilenceUsage:  true,
	SilenceErrors: true,
	// 所有操作都需要登录，所以放在这里
	// PreRun 达不到效果
	// 这里的代码会被子命令覆盖，要注意！！！
	//PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
	//},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// 设定默认值，在项目根目录下执行可以
	// 直接在CMD窗口执行，无法获取该配置文件
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file (default is $HOME/.goeds.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goeds.yaml)")

	rootCmd.PersistentFlags().StringP("logger-file", "f", "", "EDS log content file")

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("EDS")

	viper.BindPFlag("logger-file", rootCmd.PersistentFlags().Lookup("logger-file"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".goeds" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".goeds")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
