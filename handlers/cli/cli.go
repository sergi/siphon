package cli

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/sergi/siphon"
)

const pingPeriod = 5 * time.Minute

type Handler struct {
	mu     sync.Mutex
	Writer io.Writer
}

func New(options ...func(*Handler)) func(*siphon.Server) {
	return func(server *siphon.Server) {
		h := &Handler{
			Writer: os.Stdout,
		}

		for _, option := range options {
			option(h)
		}

		go h.handleConnection(server)
	}
}

func (h *Handler) handleConnection(s *siphon.Server) {
	sub := s.Subscribe(siphon.DefaultChannel)
	t := time.NewTicker(pingPeriod)

	var message []byte

	for {
		select {
		case <-t.C:
			message = nil
		case message = <-sub:
		}

		_, err := h.Writer.Write(message)
		if err != nil {
			fmt.Print(err)
			break
		}
	}

	t.Stop()
	s.Unsubscribe(siphon.DefaultChannel, sub)
}
