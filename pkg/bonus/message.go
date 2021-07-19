package bonus

import (
	"github.com/golang/protobuf/proto"
	"github.com/newswarm-lab/new-bee/pkg/bonus/message"

)

type message_ struct {
	id  message.CSID
	msg proto.Message
}

type writeMsg struct {
	msg   *message_
	errCh chan error
}

type readMsg struct {
	msg *message_
	err error
}

type writeCh chan *writeMsg
type readCh chan *readMsg

type closeCh chan struct{}
