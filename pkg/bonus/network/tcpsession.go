package network

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/newswarm-lab/new-bee/pkg/bonus/log"
)

var (
	curSessionID = uint32(0)
	sessionLock  sync.Mutex
)

type PingPongFunc func(session *Session)

type Session struct {
	sync.Mutex
	waitGroup        sync.WaitGroup //sync
	sessionID        uint32         //sessionid
	conn             net.Conn
	writePending     chan []byte
	pendingWriteNum  uint16
	closeCtrl        chan bool
	closeFlag        bool
	recvParser       *MsgParser
	sendParser       *MsgParser
	msgProcessor     TcpProcessor
	activeTime       int64
	Writable         bool
	pingPong         PingPongFunc
	pingPongInterval int64
}

func newSession(conn net.Conn, pendingWriteNum uint16, processor TcpProcessor) *Session {
	if pendingWriteNum <= 0 {
		pendingWriteNum = 100
		log.Info("invalid PendingWriteNum, reset to %v", pendingWriteNum)
	}
	session := &Session{
		waitGroup:       sync.WaitGroup{},
		conn:            conn,
		writePending:    make(chan []byte, pendingWriteNum),
		pendingWriteNum: pendingWriteNum,
		closeCtrl:       make(chan bool),
		closeFlag:       false,
		recvParser:      NewMsgParser(nil),
		sendParser:      NewMsgParser(nil),
		msgProcessor:    processor,
		activeTime:      time.Now().Unix(),
		Writable:        true,
	}

	sessionLock.Lock()
	curSessionID++
	if curSessionID == 0 {
		curSessionID++
	}
	session.sessionID = curSessionID
	sessionLock.Unlock()

	return session
}

func (slf *Session) SetCipher(recv, send Cipher) {
	slf.recvParser = NewMsgParser(recv)
	slf.sendParser = NewMsgParser(send)
}

func (slf *Session) Start() {
	log.Info("tcp session %d addr %s connnect succ", slf.GetID(), slf.RemoteAddr().String())
	slf.msgProcessor.OnConnectSucc(slf)
	slf.waitGroup.Add(1)
	go slf.writeLoop()
	slf.readLoop()
	slf.waitGroup.Wait()
	slf.msgProcessor.OnConnectClose(slf)
}

func (slf *Session) Close() {
	log.Info("server kickout tcp session %d", slf.GetID())
	slf.close()
}

func (slf *Session) LocalAddr() net.Addr {
	return slf.conn.LocalAddr()
}

func (slf *Session) RemoteAddr() net.Addr {
	return slf.conn.RemoteAddr()
}

func (slf *Session) GetActiveTime() int64 {
	return slf.activeTime
}

func (slf *Session) GetID() uint32 {
	return slf.sessionID
}

func (slf *Session) SetWritable(bWrite bool) {
	slf.Writable = bWrite
}

func (slf *Session) SetPingPong(pingpong PingPongFunc, interval uint32) {
	slf.pingPong = pingpong
	slf.pingPongInterval = int64(interval)
}

func (slf *Session) SyncWriteMsg(msg []byte) error {
	if !slf.Writable {
		return errors.New("session Writable false")
	}
	return slf.sendParser.Write(slf.conn, msg)
}

func (slf *Session) WriteMsg(msg []byte) error {
	msgLen := len(msg)
	if msgLen == 0 {
		log.Error("write msg len is 0")
		return nil
	}

	if !slf.Writable {
		return errors.New("session Writable false")
	}
	if uint16(len(slf.writePending)) < slf.pendingWriteNum {
		slf.writePending <- msg
		return nil
	}

	log.Error("sender overflow, pending:%d tcp session:%d", uint16(len(slf.writePending)), slf.GetID())
	return fmt.Errorf("sender overflow, tcp session:%d", slf.GetID())
}

func (slf *Session) close() {
	slf.Lock()
	defer slf.Unlock()
	if slf.closeFlag {
		return
	}
	slf.closeFlag = true
	slf.closeCtrl <- true
}

func (slf *Session) doSyncRead() ([]byte, error) {
	return slf.recvParser.Read(slf.conn)
}

func (slf *Session) writeLoop() {
	defer func() {
		log.Info("tcp session %d addr %s close write goroutine,pending:%d", slf.GetID(), slf.RemoteAddr().String(), len(slf.writePending))
		if r := recover(); r != nil {
			log.Error("writeLoop recover tcp session %d error: %+v.", slf.GetID(), r)
		}
		slf.conn.Close()
		slf.waitGroup.Done()
	}()
	if slf.pingPongInterval != 0 {
		ticker := time.NewTicker(time.Duration(slf.pingPongInterval) * time.Second)
		for {
			select {
			case <-slf.closeCtrl:
				return
			case msg := <-slf.writePending:
				slf.activeTime = time.Now().Unix()
				err := slf.SyncWriteMsg(msg)
				if err != nil {
					log.Error("tcp session %d write msg:%x error:%+v", slf.GetID(), msg, err)
					return
				}

			case <-ticker.C:
				slf.pingPong(slf)
			}
		}
	} else {
		for {
			select {
			case <-slf.closeCtrl:
				return
			case msg := <-slf.writePending:
				slf.activeTime = time.Now().Unix()
				err := slf.SyncWriteMsg(msg)
				if err != nil {
					log.Error("tcp session %d write msg:%x error:%+v", slf.GetID(), msg, err)
					return
				}
			}
		}
	}
}

func (slf *Session) readLoop() {
	defer func() {
		log.Info("tcp session %d addr %s close read goroutine", slf.GetID(), slf.RemoteAddr().String())
		if r := recover(); r != nil {
			log.Error("readLoop recover tcp session %d error: %+v.", slf.GetID(), r)
		}
		slf.conn.Close()
		slf.close()
	}()

	for {
		if slf.closeFlag {
			return
		}
		msg, err := slf.doSyncRead()
		if err != nil {
			log.Error("tcp session %d read msg eror:%+v", slf.GetID(), err)
			return
		}
		slf.activeTime = time.Now().Unix()
		slf.msgProcessor.HandleRecvMsg(slf, msg)
	}
}
