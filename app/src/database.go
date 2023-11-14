package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gopkg.in/yaml.v3"
)

type Database struct {
	Run    bool   `yaml:"running"`
	Bucket string `yaml:"bucket"`
	Org    string `yaml:"org"`
	Tocken string `yaml:"tocken"`
	Url    string `yaml:"url"`
}

type DbConfig struct {
	TestDB    Database `yaml:"db"`
	TestUnits Units    `yaml:"units"`
	lock      sync.Mutex
}

func (db *DbConfig) readDB(filename *string, config *DbConfig) error {
	log.Println("Reading config", *filename)
	buf, err := os.ReadFile(*filename)
	if err != nil {
		return err
	}
	yaml.Unmarshal(buf, &config)
	return nil
}

func (db *DbConfig) writeAPoint(wg *sync.WaitGroup, messages chan string) {
	db.lock.Lock()
	defer db.lock.Unlock()
	log.Print("Creating Client... ")
	time.Sleep(3 * time.Second)
	client := influxdb2.NewClient(db.TestDB.Url, db.TestDB.Tocken)
	writeAPI := client.WriteAPIBlocking(db.TestDB.Org, db.TestDB.Bucket)
	time.Sleep(2 * time.Second)
	fmt.Printf("Ready!\nStarting to send messages\nCheck database in %v\n", db.TestDB.Url)
	for out := range messages {
		//log.Println("whats obtained:", out)
		measurement, unit, values := splitChanOut(out)
		if values == nil {
			continue
		}
		epoch := int64(values["timestamp"])
		p := influxdb2.NewPointWithMeasurement(measurement).AddTag("unit", unit).AddField("value", values["value"]).SetTime(time.UnixMilli(epoch))
		writeAPI.WritePoint(context.Background(), p)
	}
	client.Close()
}

func (db *Database) queryDatabase() {
	client := influxdb2.NewClient(db.Url, db.Tocken)
	queryAPI := client.QueryAPI(db.Org)

	result, err := queryAPI.Query(context.Background(), `from(bucket:"networktool")|> range(start: -1h)|> filter(fn: (r) => r._measurements == "stat")`)
	if err == nil {
		for result.Next() {
			fmt.Printf("row %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query err: %s\n", result.Err().Error())
		}
	}
	client.Close()
}

func splitChanOut(out string) (string, string, map[string]float64) {
	//log.Print(out)
	var measurement string
	var values = make(map[string]float64)
	sep := strings.Split(out, ":")
	measurement = sep[0]
	//log.Print("measurement type", measurement)
	sep = strings.Split(sep[1], ",")
	timestamp, err := strconv.ParseFloat(sep[0], 64)
	//log.Print("the timestamp", timestamp)
	if err != nil {
		log.Fatal("Error parsing to Float", err)
	}
	value := sep[1]
	if len(value) != 0 {
		//log.Print("The value", value)
		sep = strings.Split(value, " ")
		new_value, err := strconv.ParseFloat(sep[0], 64)
		unit := sep[1]
		if err != nil {
			log.Fatal("Error parsing to Float", err)
		}
		values["timestamp"] = timestamp
		values["value"] = new_value
		return measurement, unit, values
	}
	log.Println("No rtt!")
	return "", "", nil
}
