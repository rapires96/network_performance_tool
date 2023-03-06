package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Process struct {
	config *Config
	channel chan string
	lock sync.Mutex
}

func (p *Process) startPing(wg *sync.WaitGroup) { //
	p.lock.Lock()
	defer p.lock.Unlock()
	//defer wg.Done()

	var command string = fmt.Sprintf("ping %s -i %v", p.config.TestPing.Server, p.config.TestPing.Interval)
	log.Printf("In Ping Scanner\nRunning: %v\n------------------\n", command)

	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panic("err initializing Pipe")
	}
	if err := cmd.Start(); err != nil {
		log.Panic("err starting the command")
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	var rtt string
	
	for scanner.Scan() {
		rtt = p.config.TestPing.getPingRTT(scanner.Text())
		p.channel <- "rtt:"+strconv.FormatInt(time.Now().UnixMilli(), 10) + "," + rtt
		time.Sleep(200*time.Millisecond)
	}

	if err := scanner.Err(); err != nil {log.Panic("Error handling stdout")}
	if err := cmd.Wait(); err != nil {log.Panic("Error waiting for command:", err)}
}

func (p *Process) startIperf(wg *sync.WaitGroup) { //
	p.lock.Lock()
	defer p.lock.Unlock()
	//defer wg.Done()

	var command string = p.config.TestIperf.makeCommand()
	log.Printf("In Iperf Scanner\nRunning: %v\n------------------\n", command)
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	
	defer cmd.Wait()
	if err != nil {	log.Panic("err initializing Pipe") }

	if err := cmd.Start(); err != nil {	log.Panic("err starting the command") }
	
	reg, _ := regexp.Compile(`(\d+.\d+) \w*bits/sec`)	

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	var bitrate string
	for scanner.Scan(){
		for _, sample := range strings.Split(scanner.Text(), "\n"){
			bitrate = p.config.TestIperf.getIperfBitrate(string(sample), reg)
		}
		if bitrate != "" {
			p.channel <- "bitrate:"+strconv.FormatInt(time.Now().UnixMilli(), 10) + "," + bitrate
		}
		time.Sleep(200*time.Millisecond)
	}

	if err = scanner.Err(); err != nil {log.Panic("error handling out ", err) }
	if err := cmd.Wait(); err != nil { log.Panic("Error waiting for command:", err) }
}