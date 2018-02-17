package siphon

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

const (
	maxMessageSize = 1024
	defaultUDPPort = 1200
	DefaultChannel = "logging"
	localhost      = "127.0.0.1"
)

// Server represents a UDP server that broadcasts all data
// chunks coming from client pipes
type Server struct {
	subscriptions      map[string][]chan []byte
	subscriptionsMutex sync.Mutex
	port               int
}

func Port(port int) func(*Server) {
	return func(server *Server) {
		if port != 0 {
			server.port = port
		}
	}
}

// NewServer creates a new server
func NewServer(options ...func(*Server)) (*Server, error) {
	s := &Server{
		port:          defaultUDPPort,
		subscriptions: map[string][]chan []byte{},
	}

	for _, option := range options {
		option(s)
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: s.port,
		IP:   net.ParseIP(localhost),
	})

	if err != nil {
		return nil, err
	}

	hostName, _ := os.Hostname()
	fmt.Println(Name + " UDP server: " + hostName + ":" + strconv.Itoa(s.port))

	s.Listen(conn)

	return s, nil
}

func (s *Server) Listen(conn *net.UDPConn) {
	for {
		var buf = make([]byte, 1500)

		n, _ /*addr*/, err := conn.ReadFromUDP(buf)
		checkError(err)
		//fmt.Println("Got message from ", addr, " with n = ", n)

		s.subscriptionsMutex.Lock()
		subs := s.subscriptions[DefaultChannel]
		s.subscriptionsMutex.Unlock()

		for _, sub := range subs {
			select {
			case sub <- buf[:n]:
			default:
				// drop the message if nobody is ready to receive it
				fmt.Println("Message dropped, nobody listening")
			}
		}
	}
}

func (s *Server) Subscribe(channel string) chan []byte {
	sub := make(chan []byte)
	s.subscriptionsMutex.Lock()
	s.subscriptions[DefaultChannel] = append(s.subscriptions[DefaultChannel], sub)
	s.subscriptionsMutex.Unlock()
	return sub
}

func (s *Server) Unsubscribe(channel string, sub chan []byte) {
	s.subscriptionsMutex.Lock()
	var newSubs []chan []byte
	subs := s.subscriptions[channel]
	for _, s := range subs {
		if s != sub {
			newSubs = append(newSubs, s)
		}
	}
	s.subscriptions[channel] = newSubs
	s.subscriptionsMutex.Unlock()
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
