package reportor

import (
	"fmt"
	"github.com/spf13/viper"
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
	viper.SetConfigName("config.yaml")
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

	return
}
