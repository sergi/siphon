package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sergi/pipe2browser"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	ExitCodeOK = iota
	ExitCodeError
)

var (
	isServer = kingpin.Flag("server", "Run in Server mode").Bool()
	port     = kingpin.Flag("port", "Port to run in Server mode").Int()
	address  = kingpin.Flag("address", "Server address").Default("127.0.0.1:3000").String()
	id       = kingpin.Flag("id", "ID given to this stream").Default("").String()
)

func main() {
	kingpin.Parse()

	if *isServer {
		server := siphon.NewServer(*port)
		err := server.Init()
		if err != nil {
			fmt.Println(err)
			os.Exit(ExitCodeError)
		}
	} else {
		err := siphon.Init(*address, *id, bufio.NewReader(os.Stdin))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(ExitCodeError)
		}
	}
	os.Exit(ExitCodeOK)
}
