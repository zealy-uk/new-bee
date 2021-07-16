package bonus

import (
	"verifycenter/log"

	"github.com/newswarm-lab/new-bee/pkg/bonus/message"
	"github.com/newswarm-lab/new-bee/pkg/bonus/network"

	"google.golang.org/protobuf/proto"
)

var (
	PeerAddr       string
	EthAddr        string
	ChequebookAddr string
)

type MyTcpProcessor struct {
	network.DefTcpProcessor
}

func (slf *MyTcpProcessor) OnConnectSucc(session *network.Session) {
	log.Info("client connect,session:%d[%s]", session.GetID(), session.RemoteAddr().String())
}
func (slf *MyTcpProcessor) OnConnectClose(session *network.Session) {
	log.Info("client close,session:%d[%s]", session.GetID(), session.RemoteAddr().String())
}

func Heartbeat(session *network.Session) {
	send := &message.Heartbeat{
		Peer:           PeerAddr,
		EthAddr:        EthAddr,
		ChequebookAddr: ChequebookAddr,
	}
	data, _ := network.PutProtobufPayload(uint16(message.CSID_ID_Heartbeat), send)
	session.WriteMsg(data)
}

func CipherKeyNtf(session *network.Session, msg proto.Message) {
	res := msg.(*message.CipherKeyNtf)
	sCipher := network.NewRc4Cipher([]byte(res.SvrKey))
	cCipher := network.NewRc4Cipher([]byte(res.CltKey))
	session.SetCipher(sCipher, cCipher)

	session.SetWritable(true)
	session.SetPingPong(Heartbeat, 15)
}

func HeartbeatRsp(session *network.Session, msg proto.Message) {
	res := msg.(*message.HeartbeatRsp)
	log.Info("recv HeartbeatRsp,%+v", res)
}
