package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sergi/siphon"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	ExitCodeOK = iota
	ExitCodeError
)

var (
	app = kingpin.New("siph", "A real-time stream utility.")

	server        = app.Command("server", "Run in Server mode")
	serverUDPPort = server.Flag("udp-port", "UDP port to run in Server mode").Short('u').Int()
	serverWsPort  = server.Flag("ws-port", "WebSockets port to run in Server mode").Short('w').Int()

	client         = app.Command("client", "Stream output to a server")
	clientAddress  = client.Flag("address", "Server address").Default("127.0.0.1:1200").String()
	clientID       = client.Flag("id", "ID given to this stream").Default("").String()
	clientNoOutput = client.Flag("no-output", "Don't emit any output. Siph re-emits its stdin output by default").Bool()
)

func main() {

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case server.FullCommand():
		server := siphon.NewServer(*serverUDPPort, *serverWsPort)
		err := server.Init()
		if err != nil {
			fmt.Println(err)
			os.Exit(ExitCodeError)
		}

	case client.FullCommand():
		err := siphon.Init(*clientAddress, *clientID, bufio.NewReader(os.Stdin), !(*clientNoOutput))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(ExitCodeError)
		}
	}

	os.Exit(ExitCodeOK)
}
