package bonus

import (
	"flag"
	"sync"

	"github.com/newswarm-lab/new-bee/pkg/bonus/log"

	"github.com/newswarm-lab/new-bee/pkg/bonus/network"
	"github.com/sirupsen/logrus"

	"github.com/newswarm-lab/new-bee/pkg/bonus/message"

	"github.com/newswarm-lab/new-bee/pkg/bonus/client"
)

func startBonus(logger *logrus.Logger) {

	var svrAddr string
	flag.StringVar(&svrAddr, "svr", "139.162.90.128:9527", "default 139.162.90.128:9527")
	flag.Parse()

	log.Init(logger)

	clientProcessor := &MyTcpProcessor{}
	clientProcessor.MsgHandles = make(map[uint16]network.MsgHander)

	clientProcessor.RegisterMsg(uint16(message.CSID_ID_CipherKeyNtf), "CipherKeyNtf", client.CipherKeyNtf)
	clientProcessor.RegisterMsg(uint16(message.CSID_ID_HeartbeatRsp), "HeartbeatRsp", client.HeartbeatRsp)

	pClient := network.NewTCPClient(svrAddr, clientProcessor)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		pClient.Start()
		wg.Done()
	}()
	log.Info("client exit")
}
