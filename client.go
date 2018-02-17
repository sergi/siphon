package siphon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/teris-io/shortid"
)

const maxBufferCapacity = 1024 * 4

type ConsumerOptions struct {
	Id         string
	Address    string
	EmitOutput bool
}

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

// Listen initializes the client to start sending Chunks at the given address
func Init(opts ConsumerOptions, stream io.Reader, conn io.Writer) error {
	if opts.Id == "" {
		opts.Id, _ = shortid.Generate()
	}

	nBytes, nChunks := int64(0), int64(0)
	if stream == nil {
		err := errors.New("no stream passed as a parameter, nowhere to read from")
		return err
	}

	host := getHostName()
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
			ID:        opts.Id,
			Data:      string(buffer),
			Host:      host,
			Timestamp: time.Now().Unix(),
		})

		if err != nil {
			break
		}

		if _, err = conn.Write(jsonToSend); err != nil && err != io.EOF {
			return err
		}

		if opts.EmitOutput == true {
			// Write to stdout
			out := os.Stdout
			if _, err = out.WriteString(string(buffer)); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	return nil
}
