package main

import (
	"log"
	"os"
	"sync"
	"time"
	"gopkg.in/yaml.v3"
)

func main() {
	log.SetOutput(os.Stdout)
	var defaultFile string =  string(os.Getenv("APP_CONFIG"))
	var configFile string
	var config Config
	if len(os.Args) == 2 { configFile = os.Args[1] } else { configFile = defaultFile }
	err := readConfig(&configFile, &config)
	messages := make(chan string, 2)
	defer close(messages)
	var wg sync.WaitGroup

	if err != nil {	log.Fatal("Error reading config") }

	go printMetrics(messages)
	
	if config.TestPing.Run {
		wg.Add(1)
		var p1 Process = Process{config: &config, channel: messages} 
		go p1.startPing(&wg)	
	}
	if config.TestIperf.Run {
		wg.Add(1)
		p2 := Process{config: &config, channel: messages}
		go p2.startIperf(&wg)
	}
	wg.Wait()
}

func printMetrics(messages chan string) { // Print metrics from the channel
	log.Println("printing the messages received from the channel ...")
	time.Sleep(2*time.Second)
	for out := range messages{
		log.Printf("%v\n", out)
	}
}

func readConfig(filename *string, config *Config) error {
	log.Println("Reading config", *filename)
	buf, err := os.ReadFile(*filename)
	if err != nil {
		return err
	}
	yaml.Unmarshal(buf, &config)
	return nil
}
