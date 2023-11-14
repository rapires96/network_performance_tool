package main

import (
	"log"
	"os"
	"sync"
	"time"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func main() {
	//err := exec.Command("bash", "./entrypoint.sh").Run()
	//if err != nil {log.Println("couldn't run entrypoint", err)}
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loadig .env file")
	}
	log.SetOutput(os.Stdout)
	var msEnv string = string(os.Getenv("APP_CONFIG")) //measurement configuration
	var dbEnv string = string(os.Getenv("DB_CONFIG"))  //database configuration
	//log.Printf("here are %v %v", msEnv, dbEnv)
	//var defaultFile string =  string(os.Getenv("APP_CONFIG"))
	var msFile, dbFile string
	var msConfig Config
	var db DbConfig
	// file s may be also passed when executinh as args first one is measurement file second one is database file
	if len(os.Args) == 2 {
		msFile = os.Args[1]
	} else {
		msFile = msEnv
	}
	if len(os.Args) == 3 {
		msFile, dbFile = os.Args[1], os.Args[2]
	} else {
		msFile, dbFile = msEnv, dbEnv
	}

	err1 := readConfig(&msFile, &msConfig)
	if err1 != nil {
		log.Fatal("Error reading measurement config", err1)
	}
	err2 := db.readDB(&dbFile, &db)
	if err2 != nil {
		log.Fatal("Error reading database config", err2)
	}

	messages := make(chan string, 2)
	defer close(messages)
	var wg sync.WaitGroup

	if msConfig.TestPing.Run {
		wg.Add(1)
		var p1 Process = Process{config: &msConfig, channel: messages}
		go p1.startPing(&wg)
	}
	if msConfig.TestIperf.Run {
		wg.Add(1)
		var p2 Process = Process{config: &msConfig, channel: messages}
		go p2.startIperf(&wg)
	}

	if db.TestDB.Run {
		wg.Add(1)
		go db.writeAPoint(&wg, messages)
	} else {
		go printMetrics(messages)
	}
	wg.Wait()
}

func printMetrics(messages chan string) { 
	// Print metrics from the channel
	log.Println("printing the messages received from the channel ...")
	time.Sleep(2 * time.Second)
	for out := range messages {
		log.Printf("%v\n", out)
	}
}

func readConfig(filename *string, config *Config) error {
	// reads the measurement configuration
	log.Println("Reading config", *filename)
	buf, err := os.ReadFile(*filename)
	if err != nil {
		return err
	}
	yaml.Unmarshal(buf, &config)
	return nil
}
