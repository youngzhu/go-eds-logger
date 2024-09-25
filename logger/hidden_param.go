package logger

import (
	"edser/http"
	"fmt"
	"log"
	"strings"
)

// 先用切片搜集数据
// 然后遍历数据统计分析
var hiddenParamObserver = make([]hiddenParamStats, 0)

type hiddenParamStats struct {
	url        string
	key        string
	value      string
	occurrence int // 出现次数
}

func (s hiddenParamStats) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb,
		"hiddenParamStats:\n\turl: %s\n\tkey: %s\n\tvalue: %s\n\toccurrence: %d\n",
		s.url, s.key, s.value, s.occurrence)

	return sb.String()
}

var hiddenParamKeys = []string{"__EVENTVALIDATION", "__VIEWSTATE"}

func getHiddenParams(url string) (map[string]string, error) {
	result := make(map[string]string)

	respHtml, err := http.DoGet(url)
	if err != nil {
		//log.Println("getHiddenParams error:", err)
		return nil, err
	}
	//println(respHtml)

	for _, k := range hiddenParamKeys {
		result[k] = getValueFromHtml(respHtml, k)
	}

	// 观察数据是否一致，该方法是否可以只调一次？
	// 每次结果都不一样！！必须重复调用！！！
	//for k, v := range result {
	//	hiddenParamObserver = append(hiddenParamObserver, hiddenParamStats{url: url, key: k, value: v, occurrence: 1})
	//}

	return result, nil
}

type hiddenParamStatsSlice struct {
	_slice []hiddenParamStats
}

func newHiddenParamStatsSlice(stats hiddenParamStats) hiddenParamStatsSlice {
	return hiddenParamStatsSlice{_slice: []hiddenParamStats{stats}}
}

// 根据Value累加
// 如果Value存在，则计数+1
// 否则增加一个对象
func (s *hiddenParamStatsSlice) addByValue(stats hiddenParamStats) {
	add := true

	for _, paramStats := range s._slice {
		if stats.value == paramStats.value {
			paramStats.occurrence = paramStats.occurrence + 1
			add = false
			break
		}
	}

	if add {
		s._slice = append(s._slice, stats)
	}

}

func (s hiddenParamStatsSlice) showStatistics() {
	// 只打印次数超过1次的数据
	for _, stats := range s._slice {
		//if stats.occurrence > 1 {
		log.Println(stats)
		//}
	}
}

func showStats() {
	// 按URL统计，看不同的URL是否有不同的值

	// 按key统计，看不同的key是否有不同的值
	groupByKey := map[string]hiddenParamStatsSlice{}
	var slice hiddenParamStatsSlice
	for _, stats := range hiddenParamObserver {
		if value, exists := groupByKey[stats.key]; exists {
			value.addByValue(stats)
			slice = value
		} else {
			slice = newHiddenParamStatsSlice(stats)
		}
		groupByKey[stats.key] = slice
	}

	for _, v := range groupByKey {
		v.showStatistics()
	}
}
