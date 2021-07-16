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
	sync.Mutex                      //同步保护
	waitGroup        sync.WaitGroup //sync
	sessionID        uint32         //sessionid
	conn             net.Conn       //网络连接
	writePending     chan []byte    //连接发送缓冲
	pendingWriteNum  uint16         //最大连接发送缓冲消息数
	closeCtrl        chan bool      //关闭控制信号
	closeFlag        bool           //连接关闭标志
	recvParser       *MsgParser     //网络消息编解码处理
	sendParser       *MsgParser     //网络消息编解码处理
	msgProcessor     TcpProcessor   //网络消息处理
	activeTime       int64          //tcp连接最近活跃时间戳
	Writable         bool           //连接可写标志
	pingPong         PingPongFunc   //pingpong心跳函数
	pingPongInterval int64          //pingpong心跳函数发送周期
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
		recvParser:      NewMsgParser(nil), //初始化为不加密的信道
		sendParser:      NewMsgParser(nil), //初始化为不加密的信道
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

//SetCipher 调用之后收发包按照新设置的加解密模块执行
func (slf *Session) SetCipher(recv, send Cipher) {
	slf.recvParser = NewMsgParser(recv)
	slf.sendParser = NewMsgParser(send)
}

//Start 开启收发包协程,
func (slf *Session) Start() {
	slf.msgProcessor.OnConnectSucc(slf)
	slf.waitGroup.Add(1)
	go slf.writeLoop()
	slf.readLoop()
	slf.waitGroup.Wait()
	slf.msgProcessor.OnConnectClose(slf)
}

//Close 服务器逻辑主动关闭连接会话
func (slf *Session) Close() {
	slf.close()
}

//LocalAddr 获取本地地址
func (slf *Session) LocalAddr() net.Addr {
	return slf.conn.LocalAddr()
}

//RemoteAddr 获取远端地址
func (slf *Session) RemoteAddr() net.Addr {
	return slf.conn.RemoteAddr()
}

//GetActiveTime 获取最近活跃时间戳
func (slf *Session) GetActiveTime() int64 {
	return slf.activeTime
}

//GetID 获取sessionid
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

//SyncWriteMsg 同步发送消息,只有刚建立连接时使用该函数，后续使用会导致乱序错误
func (slf *Session) SyncWriteMsg(msg []byte) error {
	if !slf.Writable {
		return errors.New("session Writable false")
	}
	return slf.sendParser.Write(slf.conn, msg)
}

//WriteMsg 异步发送消息,通常情况下都应该使用该接口
func (slf *Session) WriteMsg(msg []byte) error {
	if !slf.Writable {
		return errors.New("session Writable false")
	}
	if uint16(len(slf.writePending)) < slf.pendingWriteNum {
		slf.writePending <- msg
		return nil
	}

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
					log.Error("tcp session %d write msg:%x eror:%+v", slf.GetID(), msg, err)
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
					log.Error("tcp session %d write msg:%x eror:%+v", slf.GetID(), msg, err)
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
