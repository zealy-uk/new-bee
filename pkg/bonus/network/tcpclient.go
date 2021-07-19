package network

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/newswarm-lab/new-bee/pkg/bonus/log"
)

type TCPClient struct {
	sync.Mutex
	DstAddr          string
	HeartbeatExpire  int64
	PendingWriteNum  uint16
	AutoReconnect    bool
	session          *Session
	closeFlag        bool
	processor        TcpProcessor
	Pingpong         PingPongFunc
	PingpongInterval uint32
}

// NewTCPClient ...
func NewTCPClient(dstAddr string, h TcpProcessor) *TCPClient {
	pNewClient := &TCPClient{
		DstAddr:   dstAddr,
		closeFlag: false,
		processor: h,
	}
	pNewClient.init()
	return pNewClient
}

func (slf *TCPClient) Close() {
	slf.Lock()
	defer slf.Unlock()
	slf.closeFlag = true
	if slf.session != nil {
		slf.session.Close()
	}
}

func (slf *TCPClient) init() {

	if slf.HeartbeatExpire <= 0 {
		slf.HeartbeatExpire = 60
		log.Info("invalid HeartbeatExpire, reset to %v", slf.HeartbeatExpire)
	}

	if slf.PendingWriteNum <= 0 {
		slf.PendingWriteNum = 100
		log.Info("invalid PendingWriteNum, reset to %v", slf.PendingWriteNum)
	}

	if slf.processor == nil {
		log.Panic("TcpClientProcessor can not nil")
	}
}

func (slf *TCPClient) SyncWriteMsg(msg []byte) error {
	if slf.session != nil {
		return slf.session.SyncWriteMsg(msg)
	}
	return errors.New("session is nil")
}

func (slf *TCPClient) WriteMsg(msg []byte) error {
	if slf.session != nil {
		return slf.session.WriteMsg(msg)
	}
	return errors.New("session is nil")
}

func (slf *TCPClient) Start() {
reConnect:
	conn, err := net.Dial("tcp", slf.DstAddr)
	if err != nil {
		log.Error("connect tcp server error: %v", err)
		if slf.AutoReconnect && !slf.closeFlag {
			time.Sleep(1 * time.Second)
			goto reConnect
		} else {
			return
		}
	}
	slf.Lock()
	if slf.closeFlag {
		slf.Unlock()
		return
	}
	slf.session = newSession(conn, slf.PendingWriteNum, slf.processor)
	slf.session.SetPingPong(slf.Pingpong, slf.PingpongInterval)
	slf.Unlock()
	slf.session.Start()
	if slf.AutoReconnect && !slf.closeFlag {
		time.Sleep(1 * time.Second)
		goto reConnect
	}
}
