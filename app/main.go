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

type ErrorCode int16

const  (
	UnsupportedVersion ErrorCode = 35
)

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
	apiVersion := binary.BigEndian.Uint16(buff[6:8])
	var errorCode ErrorCode = 0
	if apiVersion < 0 || apiVersion > 4 {
		errorCode = UnsupportedVersion
	}
	respError := make([]byte, 2)
	binary.BigEndian.PutUint16(respError, uint16(errorCode))
	resp := make([]byte, 10)
	copy(resp, []byte{0,0,0,0})
	copy(resp[4:], buff[8:12])
	copy(resp[8:], respError)
	conn.Write(resp)
	conn.Close()
}
