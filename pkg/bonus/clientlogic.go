package bonus

import (
	"encoding/json"
	"github.com/newswarm-lab/new-bee/pkg/swarm"

	"github.com/newswarm-lab/new-bee/pkg/bonus/bonuskey"
	"github.com/newswarm-lab/new-bee/pkg/bonus/log"
	"github.com/newswarm-lab/new-bee/pkg/settlement/swap/chequebook"

	"github.com/newswarm-lab/new-bee/pkg/bonus/message"
	"github.com/newswarm-lab/new-bee/pkg/bonus/network"

	"google.golang.org/protobuf/proto"

	"github.com/newswarm-lab/new-bee/pkg/settlement/swap"
)

type MyTcpProcessor struct {
	network.DefTcpProcessor
}

func (slf *MyTcpProcessor) OnConnectSucc(session *network.Session) {
	session.SetWritable(false)
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
	data, err := network.PutProtobufPayload(uint16(message.CSID_ID_Heartbeat), send)
	if err != nil {
		log.Error("PutProtobufPayload error:%s", err.Error())
	}
	session.WriteMsg(data)
}

func (slf *MyTcpProcessor) CipherKeyNtf(session *network.Session, msg proto.Message) {
	res := msg.(*message.CipherKeyNtf)
	sCipher := network.NewRc4Cipher([]byte(res.SvrKey))
	cCipher := network.NewRc4Cipher([]byte(res.CltKey))
	session.SetCipher(sCipher, cCipher)
	session.SetWritable(true)
}

func (slf *MyTcpProcessor) HeartbeatRsp(session *network.Session, msg proto.Message) {
	_ = msg.(*message.HeartbeatRsp)
}

func (slf *MyTcpProcessor) EmitCheque(session *network.Session, msg proto.Message) {
	res := msg.(*message.EmitCheque)
	signedCheque := &chequebook.SignedCheque{}
	err := json.Unmarshal(res.Cheque, signedCheque)
	if err != nil {
		log.Error("SignedCheque Unmarshal error:%s", err.Error())
		return
	}

	peer := swarm.NewAddress(signedCheque.Chequebook.Bytes())

	if err := swap.BonusSwapService.ReceiveBonusCheque(nil, peer, signedCheque); err != nil {
		log.Error("failed to finally receive and store swap bonus cheque: chequebook:%s, chequeId:%s. ERROR: %w", signedCheque.Chequebook, signedCheque.Id, err)
	}
}
