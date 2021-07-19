package bonus

import (
	"crypto/rc4"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/newswarm-lab/new-bee/pkg/bonus/message"
)

var enEntry = encodeEntry{
	message.CSID_ID_Heartbeat: encodeHeartbeat,
}

type encoder struct {
	encipher *rc4.Cipher
	//enEntry  encodeEntry
}

func newEncoder(decipher *rc4.Cipher) *encoder {
	return &encoder{
		encipher: decipher,
		//enEntry:  enEntry,
	}
}

type encodeFn func(msg proto.Message) ([]byte, error)
type encodeEntry map[message.CSID]encodeFn

func (e *encoder) encode(msgID message.CSID, msg proto.Message, crypt bool) ([]byte, error) {
	data, err := enEntry[msgID](msg)
	if err != nil {
		return nil, err
	}
	b := make([]byte, msgIndexLen)
	binary.BigEndian.PutUint16(b, uint16(msgID))
	b = append(b, data...)

	p := b
	if crypt {
		p = make([]byte, len(b))
		e.encipher.XORKeyStream(p, b)
	}

	return p, nil
}

func encodeHeartbeat(msg proto.Message) ([]byte, error) {
	return proto.Marshal(msg)
}
