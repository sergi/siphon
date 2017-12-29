package siphon

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxMessageSize = 1024
	pingPeriod     = 5 * time.Minute
	defaultWsPort  = 3000
	defaultUDPPort = 1200
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(r *http.Request) bool { // WARNING! That's probably messed up https://stackoverflow.com/questions/33323337/update-send-the-origin-header-with-js-websockets
		return true
	},
}

// Server represents a UDP and Websockets server that broadcasts all data
// chunks coming from client pipes
type Server struct {
	subscriptions      map[string][]chan []byte
	subscriptionsMutex sync.Mutex
	connected          int64
	failed             int64
	wsPort             int
	udpPort            int
}

// NewServer creates a new server with the given port, or a. default one in
// case it's not provided.
func NewServer(udpPort int, wsPort int) *Server {
	if udpPort == 0 {
		udpPort = defaultUDPPort
	}
	if wsPort == 0 {
		wsPort = defaultWsPort
	}
	s := &Server{
		wsPort:  wsPort,
		udpPort: udpPort,
	}
	return s
}

func (s *Server) Init() error {
	s.subscriptions = map[string][]chan []byte{}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Connected client")
		// channel := r.URL.Path[1:]
		// launch a new goroutine so that this function can return and the http server can free up
		// buffers associated with this connection
		go s.handleConnection(ws, "logging")
	})

	go func() {
		if err := http.ListenAndServe(":"+strconv.Itoa(s.wsPort), nil); err != nil {
			log.Fatal(err)
		}
	}()

	s.startReceiving()
	return nil
}

func (s *Server) handleConnection(ws *websocket.Conn, channel string) {
	sub := s.subscribe(channel)
	atomic.AddInt64(&s.connected, 1)
	t := time.NewTicker(pingPeriod)

	var message []byte

	for {
		select {
		case <-t.C:
			message = nil
		case message = <-sub:
		}

		ws.SetWriteDeadline(time.Now().Add(30 * time.Second))
		err := ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	atomic.AddInt64(&s.connected, -1)
	atomic.AddInt64(&s.failed, 1)

	t.Stop()
	ws.Close()
	s.unsubscribe(channel, sub)
}

func (s *Server) startReceiving() {
	addr := net.UDPAddr{
		Port: s.udpPort,
		IP:   net.ParseIP("127.0.0.1"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	checkError(err)

	hostName, _ := os.Hostname()
	fmt.Println(Name + " WebSocket server running in " + hostName + ":" + strconv.Itoa(s.wsPort))
	fmt.Println(Name + " UDP server running in " + hostName + ":" + strconv.Itoa(s.udpPort))

	for {
		var buf = make([]byte, 1500)

		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Got message from ", addr, " with n = ", n)

		s.subscriptionsMutex.Lock()
		subs := s.subscriptions["logging"]
		s.subscriptionsMutex.Unlock()
		fmt.Println(string(buf[:n]))

		for _, s := range subs {
			select {
			case s <- buf[:n]:
			default:
				// drop the message if nobody is ready to receive it
			}
		}
	}
}

func (s *Server) subscribe(channel string) chan []byte {
	sub := make(chan []byte)
	s.subscriptionsMutex.Lock()
	s.subscriptions[channel] = append(s.subscriptions[channel], sub)
	s.subscriptionsMutex.Unlock()
	return sub
}

func (s *Server) unsubscribe(channel string, sub chan []byte) {
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
