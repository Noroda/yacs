package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"os"
	"sync"
	"time"

	"github.com/dreamscached/minequery/v2"
	"github.com/zan8in/masscan"
)

func main() {
	//context, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	//defer cancel()

	//godotenv.Load()

	//PORTRANGE := os.Getenv("PORT_RANGE")
	//IPRANGE := os.Getenv("IP_RANGE")
	IPRANGE1 := flag.String("range", "127.0.0.1", "IP range to scan")
	PORTRANGE1 := flag.String("port-range", "25565", "Port range to scan")
	OUTFILE1 := flag.String("output", "output.txt", "You can't disable it")
	RATE1 := flag.Int("rate", 1000, "masscan rate")
	flag.Parse()
	IPRANGE2 := *IPRANGE1
	PORTRANGE2 := *PORTRANGE1
	OUTFILE2 := *OUTFILE1
	RATE2 := *RATE1

	var (
		scannerResult []masscan.ScannerResult
		errorBytes    []byte
	)

	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets(IPRANGE2),
		masscan.SetParamPorts(PORTRANGE2),
		masscan.EnableDebug(),
		masscan.SetParamWait(0),
		masscan.SetParamRate(RATE2),
		masscan.SetParamExclude("255.255.255.255"),
	)

	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}

	if err := scanner.RunAsync(); err != nil {
		fmt.Println(err)
	}

	stdout := scanner.GetStdout()

	stderr := scanner.GetStderr()

	//fuckYou := func(text chat.Message) string {
	//      if strings.ContainsAny(fmt.Sprint(text), "\"") == true {
	//              text := strings.ReplaceAll(fmt.Sprint(text), "\"", "'")
	//              return text
	//      }
	//      return fmt.Sprint(text)
	//}

	scanAndInsert := func(ip string, port string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered!")
			}
		}()
		pinger := minequery.NewPinger(
			minequery.WithTimeout(5*time.Second),
			minequery.WithUseStrict(true),
			minequery.WithProtocolVersion16(minequery.Ping16ProtocolVersion152),
			minequery.WithProtocolVersion17(minequery.Ping17ProtocolVersion172),
		)
		portint, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("port conversion gave an error:" + err.Error())
		}
		status, err := pinger.Ping16(ip, portint)
		if status != nil {
			fmt.Println("IP: " + ip + ":" + port)
			fmt.Println("MOTD: " + status.MOTD)
			f, err := os.OpenFile(fmt.Sprint(OUTFILE2), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err == nil {
				f.WriteString("IP: " + ip + ":" + port + "\n" + "MOTD: " + status.MOTD + "\n")
			}
			var wg sync.WaitGroup
			if err != nil {
				fmt.Println(err)
			}
			wg.Wait()
		}
		if err != nil {
			fmt.Println(err)
		}
	}

	go func() {
		for stdout.Scan() {
			srs := masscan.ParseResult(stdout.Bytes())
			scannerResult = append(scannerResult, srs)
			go scanAndInsert(srs.IP, srs.Port)
			//f, err := os.OpenFile(fmt.Sprint(OUTFILE2), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			//if err == nil {
			//f.WriteString("IP: " + srs.IP + ":" + srs.Port + fmt.Sprint(&s))
			//}
		}
	}()

	go func() {
		for stderr.Scan() {
			fmt.Println(stderr.Text())
			errorBytes = append(errorBytes, stderr.Bytes()...)
		}
	}()

	if err := scanner.Wait(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("masscan result count : ", len(scannerResult))
}
