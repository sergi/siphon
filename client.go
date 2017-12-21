package p2b

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
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

func getIdentifier() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = getOutboundIP().String()
	}
	//ID, _ := shortid.Generate()
	return hostname //+ "-" + ID
}

var sessionID = getIdentifier()

func Init(address string) error {
	fmt.Println("connecting to " + address)
	udpAddr, err := net.ResolveUDPAddr("udp4", address)

	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		return err
	}

	nBytes, nChunks := int64(0), int64(0)
	streamReader := bufio.NewReader(os.Stdin)

	buffer := make([]byte, 0, maxBufferCapacity)
	for {
		n, err := streamReader.Read(buffer[:cap(buffer)])
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
			//ID:        guid,
			Data:      string(buffer),
			Host:      sessionID,
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

		// Write to stdout
		out := os.Stdout
		if _, err = out.WriteString(string(buffer)); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	return nil
}
