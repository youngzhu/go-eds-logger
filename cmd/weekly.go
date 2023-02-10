/*
Copyright © 2023 youngzy
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// weeklyCmd represents the weekly command
var weeklyCmd = &cobra.Command{
	Use:   "weekly",
	Short: "填写这一周的日志",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("weekly called")
	},
}

func init() {
	rootCmd.AddCommand(weeklyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// weeklyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// weeklyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
