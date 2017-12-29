package siphon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/teris-io/shortid"
)

const maxBufferCapacity = 1024 * 4

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func getHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = getOutboundIP().String()
	}
	return hostname
}

func getUDPAddress(address string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", address)

	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Init initializes the client to start sending Chunks at the given address
func Init(address string, id string, stream *bufio.Reader, emitOutput bool) error {
	if id == "" {
		id, _ = shortid.Generate()
	}

	conn, err := getUDPAddress(address)
	if err != nil {
		return err
	}

	nBytes, nChunks := int64(0), int64(0)
	if stream == nil {
		stream = bufio.NewReader(os.Stdin)
	}

	buffer := make([]byte, 0, maxBufferCapacity)
	for {
		n, err := stream.Read(buffer[:cap(buffer)])
		buffer = buffer[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		nChunks++
		nBytes += int64(len(buffer))

		jsonToSend, err := json.Marshal(Chunk{
			ID:        id,
			Data:      string(buffer),
			Host:      getHostName(),
			Timestamp: time.Now().Unix(),
		})

		if err != nil {
			break
		}

		_, err = conn.Write(jsonToSend)

		if err != nil {
			return err
		}
		// process buffer
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if emitOutput == true {
			// Write to stdout
			out := os.Stdout
			if _, err = out.WriteString(string(buffer)); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	return nil
}
