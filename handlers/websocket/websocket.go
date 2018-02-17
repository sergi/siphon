package websocket

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sergi/siphon"
)

const defaultWsPort = 3000

const (
	maxMessageSize = 1024
	pingPeriod     = 5 * time.Minute
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(r *http.Request) bool { // WARNING! That's probably messed up https://stackoverflow.com/questions/33323337/update-send-the-origin-header-with-js-websockets
		return true
	},
}

type Handler struct {
	connected int64
	failed    int64
	wsPort    int
}

func New(options ...func(*Handler)) func(*siphon.Server) {
	return func(server *siphon.Server) {
		h := &Handler{
			wsPort: defaultWsPort,
		}

		for _, option := range options {
			option(h)
		}

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			ws, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			// channel := r.URL.Path[1:]
			// launch a new goroutine so that this function can return and the http server can free up
			// buffers associated with this connection
			go h.handleConnection(ws, server)
		})

		go func() {
			log.Println("Initializing Websocket Handler at localhost:" + strconv.Itoa(h.wsPort) + "/ws")
			if err := http.ListenAndServe(":"+strconv.Itoa(h.wsPort), nil); err != nil {
				log.Fatal(err)
			}
		}()
	}
}

func Port(port int) func(*Handler) error {
	return func(h *Handler) error {
		h.wsPort = port
		return nil
	}
}

func (h *Handler) handleConnection(ws *websocket.Conn, s *siphon.Server) {
	atomic.AddInt64(&h.connected, 1)
	sub := s.Subscribe(siphon.DefaultChannel)
	t := time.NewTicker(pingPeriod)

	var message []byte

	for {
		select {
		case <-t.C:
			message = nil
		case message = <-sub:
		}

		ws.SetWriteDeadline(time.Now().Add(30 * time.Second))
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
	atomic.AddInt64(&h.connected, -1)
	atomic.AddInt64(&h.failed, 1)

	t.Stop()
	ws.Close()
	s.Unsubscribe(siphon.DefaultChannel, sub)
}
