package bonus

import (
	"sync"

	"github.com/newswarm-lab/new-bee/pkg/bonus/log"
	"github.com/newswarm-lab/new-bee/pkg/logging"

	"github.com/newswarm-lab/new-bee/pkg/bonus/network"

	msg "github.com/newswarm-lab/new-bee/pkg/bonus/message"
)

type BonusClient struct {
	wg sync.WaitGroup
	c  *network.TCPClient
}

func InitBonus(logger logging.Logger) *BonusClient {
	bonusClient := &BonusClient{}

	log.Init(logger)

	clientProcessor := &MyTcpProcessor{}
	clientProcessor.MsgHandles = make(map[uint16]network.MsgHander)

	clientProcessor.RegisterMsg(uint16(msg.CSID_ID_CipherKeyNtf), "CipherKeyNtf", clientProcessor.CipherKeyNtf)
	clientProcessor.RegisterMsg(uint16(msg.CSID_ID_HeartbeatRsp), "HeartbeatRsp", clientProcessor.HeartbeatRsp)
	clientProcessor.RegisterMsg(uint16(msg.CSID_ID_EmitCheque), "EmitCheque", clientProcessor.EmitCheque)

	bonusClient.c = network.NewTCPClient("", clientProcessor)
	bonusClient.c.Pingpong = clientProcessor.Heartbeat
	bonusClient.c.PingpongInterval = 15
	bonusClient.c.AutoReconnect = true

	bonusClient.wg.Add(1)
	go func() {
		defer bonusClient.wg.Done()
		bonusClient.c.Start()
	}()

	return bonusClient
}

func (bc *BonusClient) Close() error {
	bc.wg.Wait()
	return nil
}
