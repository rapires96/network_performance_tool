package main

import (
	//"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
)


type Metric struct{
	Epoch int `json:"timestamp"`
	Value map[string]int `json:"value"`
	Units string `json:"unit"`
}

type Ping struct {
	Run bool `yaml:"running"`
	Server string `yaml:"server"`
	Interval float32 `yaml:"interval"`
}

type Iperf struct {
	Run bool `yaml:"running"`
	Server string `yaml:"server"`
	Interval float32 `yaml:"interval"`
	Direction string `yaml:"direction"`
	Bitrate string `yaml:"bitrate"`
	Duration string `yaml:"time"`
	Traffic string `yaml:"traffic"`
	BufferLen string `yaml:"buffer_len"`
	Window string `yaml:"window"`
	Port int `yaml:"port"`
	Parallel int `yaml:"parallel"`
}

type Units struct{
	RTT string 	`yaml:"rtt"`
	Throughput string `yaml:"throughput"`
}

type Config struct {
	TestPing Ping `yaml:"ping"`
	TestIperf Iperf `yaml:"iperf"`
	TestUnits Units `yaml:"units"`
}

func (iperf *Iperf) makeCommand() string{
	// crafts commands for the iperf process
	if iperf.Traffic == "udp" && iperf.Direction == "uplink" {
		return fmt.Sprintf("iperf3 -u -f M -c %s -p %d -i %v -t %v -b %s -l %s --forceflush", iperf.Server, iperf.Port, iperf.Interval, iperf.Duration, iperf.Bitrate, iperf.BufferLen)
	} else if iperf.Traffic == "udp" && iperf.Direction == "downlink"{ // "tcp"
		return fmt.Sprintf("iperf3 -u -f M -R -c %s -p %d -i %v -t %v -b %s -l %s --forceflush", iperf.Server, iperf.Port, iperf.Interval, iperf.Duration, iperf.Bitrate, iperf.BufferLen)
	} else if iperf.Traffic == "tcp" && iperf.Direction == "uplink" {
		return fmt.Sprintf("iperf3 -f M -c %s -p %d -i %v -t %v -P %d -w %s --forceflush", iperf.Server, iperf.Port, iperf.Interval, iperf.Duration, iperf.Parallel, iperf.Window)
	} else if iperf.Traffic == "tcp" && iperf.Direction == "downlink" {
		return fmt.Sprintf("iperf3 -f M -R -c %s -p %d -i %v -t %v -P %d -w %s --forceflush", iperf.Server, iperf.Port, iperf.Interval, iperf.Duration, iperf.Parallel, iperf.Window)
	} else {
		log.Panic("Typo in the iperf config yaml files")
	}
	return ""
}

func (ping *Ping) getPingRTT(out string) string{
	//filters out the RTT for the ping command
	//var s string = string(out)
	var result string
	lines := strings.Split(out, "\n")
	//var targetLine string
	regc, err := regexp.Compile(`time=(\d+.\d+) ms`)
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range lines{
		if regc.MatchString(line){
			matches := regc.FindAllString(line, -1)
			for _, match := range matches{
				result = strings.Split(match, "=")[1]
			}
		}
	}
	return result
}

func (iperf *Iperf) getIperfBitrate(data string, reg *regexp.Regexp) string {
	var result string
	var line string = string(data)
	//reg, _ := regexp.Compile(`(\d+.\d+) \w*bits/sec`)
	if iperf.Parallel > 1 {
		if strings.Contains(line, "SUM"){ result = reg.FindString(line)}
	} else {
		result = reg.FindString(line)
	}
	return result
}