package speedtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

func startClient(testType, host, port string) error {
	serverAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	switch testType {
	case "download":
		conn.SetReadBuffer(LenBufJSON + LenBufData)
		downloadClient(conn)
	case "upload":
	}
	return nil
}

func readHeader(conn *net.TCPConn, buffer []byte) (*Header, error) {
	buffer = buffer[:LenBufJSON]
	_, err := io.ReadFull(conn, buffer)
	if err != nil {
		fmt.Println("read error")
		fmt.Println(err)
		return nil, err
	}
	buffer = bytes.TrimRight(buffer, "\x00")

	var header *Header
	err = json.Unmarshal(buffer, &header)
	if err != nil {
		fmt.Println("json unmarshal error")
		fmt.Println(err)
	}
	return header, err
}

func readData(conn *net.TCPConn, buffer []byte) error {
	buffer = buffer[:LenBufData]
	//fmt.Println(len(buffer))
	_, err := io.ReadFull(conn, buffer)
	if err != nil {
		fmt.Println("read error")
		fmt.Println(err)
		return err
	}

	//fmt.Println("bytes read", len(buffer))
	return nil
}

func downloadClient(conn *net.TCPConn) {
	fmt.Println("starting download speed test")
	conn.Write([]byte("download"))

	bufferData := make([]byte, LenBufData)
	lastHeard := time.Now()
	var records []Record
	var downloadBegin time.Time

	defer func() { //fmt.Printf("records: %#v\n", records)
		for _, record := range records {
			fmt.Printf("Time (ms) (%d),\tSize: %d\n", record.TimeSlot.Milliseconds(), record.Size)
		}
		fmt.Println("ending download speed test")
		analyzeResults(records)
	}()

	for {
		header, err := readHeader(conn, bufferData)
		if err != nil {
			if err == io.EOF {
				return
			}
			if time.Since(lastHeard) > time.Second*5 {
				fmt.Println("time out")
				fmt.Println(err)
				return
			}
			continue
		}

		lastHeard = time.Now()
		switch header.Type {
		case Start:
			downloadBegin = time.Now()
			records = append(records, Record{TimeSlot: time.Since(downloadBegin), Size: LenBufJSON})
		case End:
			records = append(records, Record{TimeSlot: time.Since(downloadBegin), Size: LenBufJSON})
			return
		case Data:
			if err = readData(conn, bufferData); err != nil {
				fmt.Println("read data error")
				fmt.Println(err)
				continue
			}
			records = append(records, Record{TimeSlot: time.Since(downloadBegin), Size: LenBufJSON + LenBufData})
		default:
			fmt.Println("other")
		}
	}

}

func analyzeResults(records []Record) {

	// calculate throughput
	var sum int32
	for _, record := range records {
		sum = sum + record.Size
	}
	totalTime := records[len(records)-1].TimeSlot.Seconds()
	mbRecieved := float64(sum) / float64(1000000)
	fmt.Printf("\trecieved %.4f mb  in %.2f seconds\n", mbRecieved, totalTime)
	fmt.Printf("\tdownload speed: %.4f mb/s\n", mbRecieved/totalTime)

}
