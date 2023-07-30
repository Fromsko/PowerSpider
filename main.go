package main

import (
	"PowerSpider/config"
	"PowerSpider/core"
	"log"
)

func main() {
	conf := config.Config{
		User: "202127530334",
		Pwd:  "102018",
		Timer: config.Timer{
			TimeUnit: "hourse",
			TimeInfo: 2,
		},
		ResDir:  "res",
		Porxy:   "http://localhost:7980",
		BaseUrl: "http://10.13.14.20:9999/",
	}
	config.InitConfig(&conf)

	// Start appliction
	if err := core.Start(); err != nil {
		log.Fatalf("[Error] %s", err)
	}
}
