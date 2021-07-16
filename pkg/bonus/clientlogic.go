package bonus

import (
	"encoding/json"

	"github.com/newswarm-lab/new-bee/pkg/bonus/bonuskey"
	"github.com/newswarm-lab/new-bee/pkg/bonus/log"
	"github.com/newswarm-lab/new-bee/pkg/settlement/swap/chequebook"

	"github.com/newswarm-lab/new-bee/pkg/bonus/message"
	"github.com/newswarm-lab/new-bee/pkg/bonus/network"

	"google.golang.org/protobuf/proto"
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

func (slf *MyTcpProcessor) Heartbeat(session *network.Session) {
	send := &message.Heartbeat{
		Peer:           bonuskey.PeerAddr,
		EthAddr:        bonuskey.EthAddr,
		ChequebookAddr: bonuskey.ChequebookAddr,
	}
	data, _ := network.PutProtobufPayload(uint16(message.CSID_ID_Heartbeat), send)
	session.WriteMsg(data)
}

func (slf *MyTcpProcessor) CipherKeyNtf(session *network.Session, msg proto.Message) {
	res := msg.(*message.CipherKeyNtf)
	sCipher := network.NewRc4Cipher([]byte(res.SvrKey))
	cCipher := network.NewRc4Cipher([]byte(res.CltKey))
	session.SetCipher(sCipher, cCipher)

	session.SetWritable(true)
	session.SetPingPong(slf.Heartbeat, 15)
}

func (slf *MyTcpProcessor) HeartbeatRsp(session *network.Session, msg proto.Message) {
	res := msg.(*message.HeartbeatRsp)
	log.Info("recv HeartbeatRsp,%+v", res)
}

func (slf *MyTcpProcessor) EmitCheque(session *network.Session, msg proto.Message) {
	res := msg.(*message.EmitCheque)
	//log.Info("%+v", res)

	signedCheque := &chequebook.SignedCheque{}
	err := json.Unmarshal(res.Cheque, signedCheque)
	if err != nil {
		log.Error("SignedCheque Unmarshal error:%s", err.Error())
		return
	}
	log.Info("recv SignedCheque,%+v", signedCheque)
}
