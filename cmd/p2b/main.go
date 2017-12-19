package main

import (
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
	address  = kingpin.Flag("address", "Server address").String()
)

func main() {
	kingpin.Parse()

	if *isServer {
		server := p2b.NewServer(*port)
		err := server.Init()
		if err != nil {
			fmt.Println(err)
			os.Exit(ExitCodeError)
		}
		os.Exit(ExitCodeOK)
	}

	p2b.Init(*address)
}
