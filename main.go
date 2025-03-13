package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/pule1234/VideoForge/config"
	"github.com/spf13/viper"
	"log"
)

var GlobalConfig *config.Config

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	//自动检查环境
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error init reading config file, %s", err)
	}
	err = viper.Unmarshal(&GlobalConfig)
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("httpbin.org"),
		colly.MaxDepth(2),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		//c.Visit(link)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit("http://httpbin.org/links/20/3")
}
