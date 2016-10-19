package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type SocketData struct {
	PackageName string
	PackageHash string
	DeviceID string
	Message  string
}

var Clients map[string]net.Conn // TODO: change to map package hash -> conn, deviceID
var connections []net.Conn

func ListenForClients(ln net.Listener) {
	messages := make(chan string, 100)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Closing connection due to:", err)
			conn.Close()
		}
		go handleConnection(conn, messages)
	}
}

// TODO: limit stream
func handleConnection(conn net.Conn, messages chan string) {
	r := bufio.NewReader(conn)
STOP:
	for {
		b, err := r.ReadBytes('\n')
		if err != nil {
			break STOP
		}
		sd := &SocketData{}
		err = json.Unmarshal(b, sd)
		if err != nil {
			fmt.Println(err)
			break STOP
		}

		go func(packageHash string, conn net.Conn) {
			Clients[packageHash] = conn // Add device connection to map
		}(sd.PackageHash, conn)
	}
	conn.Close()
}

// TODO: Add broadcastMessage ability
// func broadcastMessage(msgs chan string) {
// 	for {
// 		select {
// 		case message := <-msgs:
// 			fmt.Println("Broadcasting", message)
// 			message += "\n"
// 			for id, connection := range clients {
// 				_, err := connection.Write([]byte(message))
// 				if err != nil {
// 					fmt.Println("Removing client due to connection error")
// 					connection.Close()
// 					delete(clients, id)
// 				}
// 				fmt.Println("Sent to " + id);
// 			}

// 		}
// 	}
// }