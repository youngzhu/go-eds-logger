module goeds

go 1.16

require (
	github.com/PuerkitoBio/goquery v1.7.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.15.0
	github.com/youngzhu/go-smail v0.1.2
	github.com/youngzhu/godate v0.7.6
)

// 网络不好时用
replace github.com/youngzhu/godate => ../godate
