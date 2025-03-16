package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type ApiRequestMessage struct {
	Length        int32
	ApiKey        int16
	ApiVersion    int16
	CorrelationId int32
	ClientIdLen   int16
	ClientId      string
}

type APIVersion struct {
	ApiKey     int16
	MinVersion int16
	MaxVersion int16
}

type ApiVersionResponse struct {
	Length           int32
	CorrelationID    int32
	ErrorCode        int16
	APIVersionsCount int16
	APIVersions      []APIVersion
}

type ErrorCode int16

const (
	NoError            ErrorCode = 0
	UnsupportedVersion ErrorCode = 35
)

func Read(r io.Reader) (*ApiRequestMessage, error) {
	var req ApiRequestMessage

	err := binary.Read(r, binary.BigEndian, &req.Length)
	if err != nil {
		fmt.Println("Error reading length: ", err.Error())
		return nil, err
	}

	expectedLength := req.Length

	err = binary.Read(r, binary.BigEndian, &req.ApiKey)
	if err != nil {
		fmt.Println("Error reading api key: ", err.Error())
		return nil, err
	}
	expectedLength -= 2

	err = binary.Read(r, binary.BigEndian, &req.ApiVersion)
	if err != nil {
		fmt.Println("Error reading api version: ", err.Error())
		return nil, err
	}
	expectedLength -= 2

	err = binary.Read(r, binary.BigEndian, &req.CorrelationId)
	if err != nil {
		fmt.Println("Error reading correlation id: ", err.Error())
		return nil, err
	}
	expectedLength -= 4

	err = binary.Read(r, binary.BigEndian, &req.ClientIdLen)
	if err != nil {
		fmt.Println("Error reading correlation id: ", err.Error())
		return nil, err
	}
	expectedLength -= 2

	clientId, err := io.ReadAll(r)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading dynamic bytes:", err)
		return nil, err
	}
	req.ClientId = string(clientId)
	fmt.Println(req)

	// Count the bytes and cast to int32
	expectedLength -= int32(len(clientId))

	if expectedLength != 0 {
		return nil, fmt.Errorf("length mismatch: expected %d more bytes", expectedLength)
	}

	return &req, nil
}

func handleClient(c net.Conn) {
	req, err := Read(c)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}
	c.Write([]byte{byte(req.CorrelationId)})
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	handleClient(conn)

	// buff := make([]byte, 1024)
	// conn.Read(buff)
	// correlationId := binary.BigEndian.Uint32(buff[8:12])
	// apiVersion := binary.BigEndian.Uint16(buff[6:8])
	// fmt.Printf("Received message %d", correlationId)

	// var errorCode ErrorCode = NoError
	// if apiVersion > 4 {
	// 	errorCode = UnsupportedVersion
	// }

	// respError := make([]byte, 2)
	// binary.BigEndian.PutUint16(respError, uint16(errorCode))
	// resp := make([]byte, 10)
	// copy(resp, []byte{0,0,0,0})
	// copy(resp[4:], buff[8:12])
	// copy(resp[8:], respError)
	// conn.Write(resp)
}
