package bonus

import (
	"net"
)

type session struct {
	*parser
	net.Conn
}

func newSession(c net.Conn) (*session, error) {
	p, err := newParser(c)
	if err != nil {
		return nil, err
	}

	return &session{
		parser: p,
		Conn:   c,
	}, nil
}
