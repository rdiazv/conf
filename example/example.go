package main

import (
  "github.com/rdiazv/conf"
  "fmt"
)

var Config = struct {
  Name string `default:"Test"`
  Email string `required:"true"`
}{}

func main() {
  conf.Load(&Config, "./conf.json")

  fmt.Println(Config)
}
