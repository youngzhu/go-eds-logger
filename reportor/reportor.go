package reportor

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/youngzhu/godate"
	"github.com/youngzhu/godate/chinese"
	"log"
	"os"
	"strings"
)

func init() {
	readConfigViaViper()
}

// 通过 viper 读取配置
func readConfigViaViper() {
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("EDS")

	viper.AutomaticEnv() // read in environment variables that match

	viper.AddConfigPath(".")  // 从 main 运行
	viper.AddConfigPath("..") // 执行单元测试

	// 配置文件的默认名：config
	// 默认首先使用 json ，所以这里用完整的文件名
	viper.SetConfigName("config.yml")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理错误
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到
			fmt.Fprintln(os.Stderr, err)
		} else {
			// 配置文件找到但解析错误
		}
	}
	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
}

// Run 填写工作周报和日报
func Run() (err error) {

	// 登录
	err = login()
	if err != nil {
		return
	}

	// 获取周报内容
	// 可与登录同步进行
	loadWorkReport()

	// 填周报
	// 还是要取当周的工作日，因为不一定都在周一执行，如服务器故障等
	today := godate.Today()

	monday := today.Workdays()[0]
	err = fillWeeklyReport(monday.String())
	if err != nil {
		return
	}

	// 填日报
	// 直接填7天日报
	for i := 0; i < 7; i++ {
		date, _ := monday.AddDay(i)
		if chinese.IsWorkDayInChina(date) {
			fillDailyReport(date.String())
		} else {
			log.Println(date, "放假")
		}
	}

	return
}
