package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/zn11ch/g001/internal/app"
	"io/ioutil"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./config/production.json", "path to config file")
}

func main() {

	config := app.NewConfig()

	configFile, _ := ioutil.ReadFile(configPath)
	err := json.Unmarshal([]byte(configFile), &config)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	done := make(chan bool)
	lines := app.NewLinesChannel()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, application := range config.Applications {
		switch application.Type {
		case "consumer":
			for i := 0; i < application.Workers; i++ {
				consumer, err := app.NewConsumer(config, lines, &application, i)
				if err != nil {
					log.Fatal(err)
				}
				defer consumer.Close()
			}
		case "producer":
			err = watcher.Add(application.Dir)
			if err != nil {
				log.Fatal(err)
			}
			for i := 0; i < application.Workers; i++ {
				producer, err := app.NewProducer(config, watcher, lines, &application, i)
				if err != nil {
					log.Fatal(err)
				}
				defer producer.Close()
			}
		}
	}

	<-done

}
