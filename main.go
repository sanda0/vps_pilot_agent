package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/sanda0/vps_pilot_agent/dto"
	"github.com/sanda0/vps_pilot_agent/tcp_client"

	"github.com/sanda0/vps_pilot_agent/services"
)

func main() {
	host := flag.String("h", "127.0.0.1", "host")
	port := flag.Int("p", 55001, "port")
	interval := flag.Int("i", 5, "interval")
	backoffSeconds := flag.Int("b", 10, "backoff seconds")
	flag.Parse()

	config := dto.Config{
		Host:     *host,
		Port:     *port,
		Interval: *interval,
	}

	var conn net.Conn
	var err error
	msgChan := make(chan dto.Msg, 100)
	reconnectChan := make(chan struct{}, 2)
	backoff := time.Duration(*backoffSeconds) * time.Second
	var cancelFunc context.CancelFunc

	for {

		conn, err = tcp_client.ConnectToTCPServer(config.Host, config.Port)
		if err != nil {
			fmt.Printf("Error connecting to server (%s:%d): %v\n", config.Host, config.Port, err)
			fmt.Printf("Retrying in %v...\n", backoff)
			time.Sleep(backoff)
			continue
		}

		fmt.Println("Connected to server")

		var ctx context.Context
		ctx, cancelFunc = context.WithCancel(context.Background())

		go services.StartCollectSystemStat(ctx, msgChan, config.Interval)
		go tcp_client.SendMsgToTCPServer(conn, msgChan, reconnectChan)
		go tcp_client.ReadMsgFromTCPServer(conn, reconnectChan)

		<-reconnectChan
		<-reconnectChan

		fmt.Println("Disconnected from server, attempting to reconnect")
		if cancelFunc != nil {
			cancelFunc()
		}
	}

}
