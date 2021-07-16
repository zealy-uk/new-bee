package network

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"
)

var (
	MaxMsgLen = uint16(10 * 1024)
	MinMsgLen = uint16(2)
)

type MsgParser struct {
	cipher Cipher
}

func NewMsgParser(cipher Cipher) *MsgParser {
	p := new(MsgParser)
	p.cipher = cipher
	return p
}

//Read goroutine safe
func (slf *MsgParser) Read(conn net.Conn) ([]byte, error) {
	bufMsgLen := make([]byte, 2)

	conn.SetReadDeadline(time.Now().Add(time.Duration(TCP_TIMEOUT) * time.Second))

	//read len
	if _, err := io.ReadFull(conn, bufMsgLen); err != nil {
		return nil, err
	}

	//parse len
	msgLen := binary.BigEndian.Uint16(bufMsgLen)

	//check len
	if msgLen > MaxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < MinMsgLen {
		return nil, errors.New("message too short")
	}

	//read data and decrypt
	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return nil, err
	}

	if slf.cipher != nil {
		msgData = slf.cipher.Decrypt(msgData)
	}

	return msgData, nil
}

//Write goroutine safe
func (slf *MsgParser) Write(conn net.Conn, msgData []byte) error {
	msgLen := uint16(len(msgData))

	//check len
	if msgLen > MaxMsgLen {
		return errors.New("message too long")
	} else if msgLen < MinMsgLen {
		return errors.New("message too short")
	}

	//put len
	pkg := make([]byte, 2)
	binary.BigEndian.PutUint16(pkg, msgLen)

	//encrypt
	if slf.cipher != nil {
		msgData = slf.cipher.Encrypt(msgData)
	}

	//put data
	pkg = append(pkg, msgData...)

	n, err := conn.Write(pkg)

	return err
}
