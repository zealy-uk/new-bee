package bonus

import (
	"github.com/newswarm-lab/new-bee/pkg/bonus/log"
	"github.com/newswarm-lab/new-bee/pkg/logging"

	"github.com/newswarm-lab/new-bee/pkg/bonus/network"

	"github.com/newswarm-lab/new-bee/pkg/bonus/message"
)

func StartBonus(logger logging.Logger) {

	svrAddr := "139.162.90.128:9527"
	log.Init(logger)

	clientProcessor := &MyTcpProcessor{}
	clientProcessor.MsgHandles = make(map[uint16]network.MsgHander)

	clientProcessor.RegisterMsg(uint16(message.CSID_ID_CipherKeyNtf), "CipherKeyNtf", clientProcessor.CipherKeyNtf)
	clientProcessor.RegisterMsg(uint16(message.CSID_ID_HeartbeatRsp), "HeartbeatRsp", clientProcessor.HeartbeatRsp)
	clientProcessor.RegisterMsg(uint16(message.CSID_ID_EmitCheque), "EmitCheque", clientProcessor.EmitCheque)

	pClient := network.NewTCPClient(svrAddr, clientProcessor)
	pClient.Pingpong = clientProcessor.Heartbeat
	pClient.PingpongInterval = 15
	pClient.AutoReconnect = true
	go pClient.Start()
}
