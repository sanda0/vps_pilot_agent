package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/sanda0/vps_pilot_agent/dto"
	"github.com/sanda0/vps_pilot_agent/tcp_client"

	"github.com/sanda0/vps_pilot_agent/services"
)

func main() {
	host := flag.String("h", "127.0.0.1", "host")
	port := flag.Int("p", 55001, "port")
	interval := flag.Int("i", 5, "interval")
	flag.Parse()

	config := dto.Config{
		Host:     *host,
		Port:     *port,
		Interval: *interval,
	}

	fmt.Println(config)

	conn, err := tcp_client.ConnectToTCPServer(config.Host, config.Port)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	msgChan := make(chan dto.Msg, 100)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go services.CollectSystemStat(msgChan, config.Interval, wg)

	wg.Add(1)
	go tcp_client.SendMsgToTCPServer(conn, msgChan, wg)

	wg.Add(1)
	go tcp_client.ReadMsgFromTCPServer(conn, wg)

	wg.Add(1)
	wg.Wait()
}
