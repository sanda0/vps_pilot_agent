package tcp_client

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/sanda0/vps_pilot_agent/dto"
	"github.com/sanda0/vps_pilot_agent/services"
)

var canSendStats = false
var nodeID int32

func ConnectToTCPServer(host string, port int) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func SendMsgToTCPServer(conn net.Conn, msgChan chan dto.Msg, reconnect chan struct{}) {
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	sysInfo, err := services.GetSystemInfo()
	if err != nil {
		fmt.Println("Error getting system info:", err)
		reconnect <- struct{}{}
		return
	}
	sysInfoJSON, err := sysInfo.ToJSON()
	if err != nil {
		fmt.Println("Error marshalling system info:", err)
		return
	}
	err = encoder.Encode(dto.Msg{
		Msg:  "connected",
		Data: sysInfoJSON,
	})

	if err != nil {
		fmt.Println("Error encoding message:", err)
	}

	for msg := range msgChan {
		if canSendStats {
			msg.NodeId = nodeID
			err = encoder.Encode(msg)
			if err != nil {
				fmt.Println("Error encoding message:", err)
				break
			}
		}
	}

	reconnect <- struct{}{}

}

func ReadMsgFromTCPServer(conn net.Conn, reconnect chan struct{}) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)
	var msg dto.Msg
	for {
		err := decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Error decoding message:", err)
			reconnect <- struct{}{}
			return
		}
		if msg.Msg == "sys_stat" {
			canSendStats = true
			nodeID = msg.NodeId
		}
		fmt.Println("Received message:", msg.Msg)
	}
}
