package bonus

import (
	"crypto/rc4"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/newswarm-lab/new-bee/pkg/bonus/message"
)

const msgIndexLen = 2

var deEntry = decodeEntry{
	message.CSID_ID_CipherKeyNtf: decodeCipherKeyNtf,
	message.CSID_ID_HeartbeatRsp: decodeHeartbeatRsp,
	message.CSID_ID_EmitCheque:   decodeEmitCheque,
}

type decoder struct {
	decipher *rc4.Cipher
	//deEntry decodeEntry
}

type decodeFn func(int []byte) (proto.Message, error)
type decodeEntry map[message.CSID]decodeFn

func newDecoder(decipher *rc4.Cipher) *decoder {
	return &decoder{
		decipher: decipher,
		//deEntry: deEntry,
	}
}

func (d *decoder) decode(in []byte, crypt bool) (message.CSID, proto.Message, error) {
	payload := in
	if crypt {
		payload = make([]byte, len(in))
		d.decipher.XORKeyStream(payload, in)
	}

	msgID := binary.BigEndian.Uint16(payload[:msgIndexLen])
	msg, err := deEntry[message.CSID(msgID)](payload[msgIndexLen:])
	if err != nil {
		return message.CSID(msgID), nil, err
	}
	return message.CSID(msgID), msg, nil
}

func decodeCipherKeyNtf(in []byte) (proto.Message, error) {
	msg := &message.CipherKeyNtf{}
	if err := proto.Unmarshal(in, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func decodeHeartbeatRsp(in []byte) (proto.Message, error) {
	msg := &message.HeartbeatRsp{}
	if err := proto.Unmarshal(in, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func decodeEmitCheque(in []byte) (proto.Message, error) {
	msg := &message.EmitCheque{}
	if err := proto.Unmarshal(in, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
