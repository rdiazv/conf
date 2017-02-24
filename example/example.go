package main

import (
	"fmt"
	"github.com/rdiazv/conf"
)

var Config = struct {
	Debug bool

	User struct {
		Name   string `required:"true"`
		Email  string `required:"true"`
		Age    int
		Active bool `default:"true"`
	}
}{}

func main() {
	conf.Load(&Config, "./conf.json", true)

	fmt.Println(Config)
}
