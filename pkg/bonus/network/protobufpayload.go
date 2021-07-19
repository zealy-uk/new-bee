package network

import (
	"encoding/binary"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var (
	payload sync.Map
)

func RegisterPayload(msgID uint16, msgName protoreflect.FullName) error {
	if _, ok := payload.Load(msgID); ok {
		return errors.Errorf("register dup payload %d", msgID)
	}
	payload.Store(msgID, msgName)
	return nil
}

func GetProtobufPayload(data []byte) (uint16, proto.Message, error) {
	if len(data) < 2 {
		return uint16(0), nil, errors.Errorf("payload too short")
	}
	msgID := binary.BigEndian.Uint16(data[0:2])
	if msgName, ok := payload.Load(msgID); ok {
		msg, err := ParseProtobuf(msgName.(protoreflect.FullName), data[2:])
		if err != nil {
			return uint16(0), nil, err
		}
		return msgID, msg, nil
	}

	return uint16(0), nil, errors.Errorf("msgid can not find")
}

func PutProtobufPayload(msgID uint16, msg proto.Message) ([]byte, error) {
	pkg := make([]byte, 2)
	binary.BigEndian.PutUint16(pkg, msgID)
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	pkg = append(pkg, data...)
	return pkg, nil
}

func ParseProtobuf(msgName protoreflect.FullName, data []byte) (proto.Message, error) {
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if err != nil {
		return nil, err
	}
	msg := msgType.New().Interface()
	err = proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
