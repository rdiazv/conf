package main

import (
	"fmt"
	"github.com/rdiazv/conf"
)

var Config = struct {
	Name   string `default:"Test"`
	Email  string `required:"true"`
	Age    int    `default:"123"`
	Active bool   `default:"true"`
}{}

func main() {
	conf.Load(&Config, "./conf.json")

	fmt.Println(Config)
}
