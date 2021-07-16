package network

import (
	"github.com/newswarm-lab/new-bee/pkg/bonus/log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type TcpProcessor interface {
	OnConnectSucc(session *Session)
	OnConnectClose(session *Session)
	HandleRecvMsg(session *Session, msg []byte)
}

type HandleFunc func(session *Session, msg proto.Message)

type MsgHander struct {
	msgName   protoreflect.FullName
	msgHandle HandleFunc
}

type DefTcpProcessor struct {
	MsgHandles map[uint16]MsgHander
}

func (slf *DefTcpProcessor) RegisterMsg(msgID uint16, name protoreflect.FullName, handle HandleFunc) {
	id := msgID
	msghandle := MsgHander{
		msgName:   name,
		msgHandle: handle,
	}
	RegisterPayload(id, name)
	slf.MsgHandles[id] = msghandle
	log.Info("register msg :%d %s", id, name)
}

func (slf *DefTcpProcessor) HandleRecvMsg(session *Session, msgData []byte) {
	msgID, msg, err := GetProtobufPayload(msgData)
	if err != nil {
		log.Error("parse msg error:%+v", err)
		return
	}
	if handle, ok := slf.MsgHandles[msgID]; ok {
		handle.msgHandle(session, msg)
	} else {
		log.Error("session %d send invalid msgid:%d", session.GetID(), msgID)
	}
}
