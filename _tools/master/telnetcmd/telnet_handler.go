package telnetcmd

import (
	"github.com/reiver/go-telnet"
)

type sessionInput struct {
	Session *session
	Str     string
}

type Handler struct {
	Input chan *sessionInput
}

func NewHandler() *Handler {
	return &Handler{
		Input: make(chan *sessionInput, 9),
	}
}

// ServeTELNET 需要考虑session
func (h *Handler) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	newSession().handle(h, ctx, w, r)
}

func (h *Handler) pushInput(s *session, str string) {
	data := &sessionInput{
		Session: s,
		Str:     str,
	}
	h.Input <- data
}
