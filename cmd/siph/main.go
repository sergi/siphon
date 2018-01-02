package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/sergi/siphon"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	exitCodeOK = iota
	ExitCodeError
)

const siphonTemplate = `
	 _____  _         _
	/  ___|(_)       | |
	\ ` + "`" + `--.  _  _ __  | |__    ___   _ __
	` + " `" + `--. \| || '_ \ | '_ \  / _ \ | '_ \
	/\__/ /| || |_) || | | || (_) || | | |
	\____/ |_|| .__/ |_| |_| \___/ |_| |_|
	          | |
	          |_|
			  
`

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

func configure() {
	buffer := bytes.NewBufferString("")

	fmt.Fprintf(buffer, siphonTemplate)
	fmt.Fprintf(buffer, kingpin.DefaultUsageTemplate)

	app.UsageTemplate(buffer.String()).Version("0.1").Author("Sergi Mansilla")
}

func main() {
	configure()
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case server.FullCommand():
		server := siphon.NewServer(*serverUDPPort, *serverWsPort)
		fmt.Print(siphonTemplate)
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

	os.Exit(exitCodeOK)
}
