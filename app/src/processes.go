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
	defer cmd.Wait()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	var rtt string
	
	for scanner.Scan() {
		rtt = p.config.TestPing.getPingRTT(scanner.Text())
		//log.Println("from ping:", rtt)
		if rtt == "" {continue}
		p.channel <- "rtt:"+strconv.FormatInt(time.Now().UnixMilli(), 10) + "," + rtt
		time.Sleep(200*time.Millisecond)
	}

	if err := scanner.Err(); err != nil {log.Panic("Error handling stdout: ping pipe\n", err)}
	if err := cmd.Wait(); err != nil {log.Panic("Error waiting for command: ping\n Verify your internet connectivity\n", err)}
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
	
	reg, _ := regexp.Compile(`(\d+\.{0,1}\d+) \w{0,1}(bits|Bytes)\/sec`)	

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	var bitrate string
	for scanner.Scan(){
		for _, sample := range strings.Split(scanner.Text(), "\n"){
			//log.Printf(sample)
			bitrate = p.config.TestIperf.getIperfBitrate(string(sample), reg)
		}
		if bitrate == ""{ continue }
		p.channel<-"throughput:"+strconv.FormatInt(time.Now().UnixMilli(), 10) + "," + bitrate
		time.Sleep(200*time.Millisecond)
		
	}

	if err = scanner.Err(); err != nil {log.Panic("error handling out: iperf3 pipe\n", err) }
	if err := cmd.Wait(); err != nil { log.Panic("Error waiting for command iperf3\n",
	"Make sure that your iperf3 server is running and it can be reache by the application\n", err) }
}

