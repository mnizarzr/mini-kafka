package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	buff := make([]byte, 1024)
	conn.Read(buff)
	fmt.Printf("Received message %v (%d)", buff[8:12], int32(binary.BigEndian.Uint32(buff[8:12])))
	resp := make([]byte, 8)
	copy(resp, []byte{0,0,0,0})
	copy(resp[4:], buff[8:12])
	conn.Write(resp)
	conn.Close()
}
