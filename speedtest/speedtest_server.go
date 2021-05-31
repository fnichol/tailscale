package speedtest

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

func StartServer(host, port string) error {
	addr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	defer l.Close()
	fmt.Println("listening on", host+":"+port, "...")

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("failed to accept")
			return err
		}

		//handle the connection in a goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) error {
	defer conn.Close()
	buf := make([]byte, 1024)
	// TODO make this use JSON
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("read error")
		return err
	}
	data := strings.TrimRight(string(buf), "\r\n\x00")

	switch data {
	case "download":
		// Start the download test
		return downloadServer(conn)
	case "upload":
	}

	return nil
}

func downloadServer(conn *net.TCPConn) error {
	fmt.Println("starting download speed test")
	startHeader := Header{Type: Start}
	// capacity that can include headers and data
	BufData := make([]byte, LenBufData, LenBufJSON+LenBufData)
	startBytes, err := marshalHeader(startHeader, false)
	if err != nil {
		return err
	}
	_, err = conn.Write(startBytes)
	if err != nil {
		return err
	}
	for startTime := time.Now(); time.Since(startTime) < downloadTestDuration; {

		// reset the slices length
		BufData = BufData[:LenBufData]
		// randomize data and get length
		lenDataGen, err := rand.Read(BufData)
		if err != nil {
			fmt.Println("fail to generate random data")
			continue
		}
		// construct and marshal header
		dataHeader := Header{Type: Data, IncomingSize: lenDataGen}
		dataBytes, err := marshalHeader(dataHeader, false)
		if err != nil {
			// do something else?
			continue
		}
		// add header in front of data
		BufData = append(dataBytes, BufData...)
		_, err = conn.Write(BufData)
		if err != nil {
			fmt.Println("error writing data to connection")
			continue
		}

	}
	endHeader := Header{Type: End}
	headerBytes, err := marshalHeader(endHeader, false)
	if err != nil {
		return err
	}
	_, err = conn.Write(headerBytes)
	if err != nil {
		return err
	}
	fmt.Println("ending download speed test")
	return nil
}

// Marshals and pads Header structs to json byte slices.
// Pads the byteslice so that its exactly LenBufJSON bytes.
func marshalHeader(header Header, debug bool) ([]byte, error) {
	b, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}
	if len(b) > LenBufJSON {
		// too big
		return nil, errors.New("too large")
	}
	padding := make([]byte, LenBufJSON-len(b))
	b = append(b, padding...)

	if debug {
		fmt.Println("sent length", len(b))
		fmt.Println("data: ", string(b))
	}

	return b, nil
}
