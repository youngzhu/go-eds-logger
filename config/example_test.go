package config_test

import (
	"fmt"
	"goeds/config"
)

var cfg config.Configuration

func init() {
	var err error
	c, err := config.Load("testdata/test.json")
	if err != nil {
		panic(err)
	}
	cfg = c
}

func ExampleDefaultConfig_GetString() {
	s, found := cfg.GetString("name")
	if found {
		fmt.Println(s)
	} else {
		fmt.Println("not found")
	}

	// Output:
	//Young
}

func ExampleDefaultConfig_GetString_notFound() {
	s, found := cfg.GetString("nickname")
	if found {
		fmt.Println(s)
	} else {
		fmt.Println("not found")
	}

	// Output:
	//not found
}

func ExampleDefaultConfig_GetString_secondary() {
	s, found := cfg.GetString("contact:mobile")
	if found {
		fmt.Println(s)
	} else {
		fmt.Println("not found")
	}

	// Output:
	//12345
}
