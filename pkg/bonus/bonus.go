package bonus

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/newswarm-lab/new-bee/pkg/settlement/swap"
	"github.com/newswarm-lab/new-bee/pkg/swarm"
	"sync"
	"time"
)

type Bonus struct {
	peerIDStr  string
	ethAdrrStr string
	beneficiaryStr string

	peer swarm.Address
	ethAdrr common.Address
	beneficiary common.Address

	closeCh closeCh
	wg      *sync.WaitGroup

	hbWriteCh    writeCh
	hbReadCh     readCh
	activeBeater *heartbeater
	//ticker *time.Ticker

	// params for establishing network connection
	addrs         map[string]string
	retryInterval time.Duration
	dialTimeout   time.Duration
	pingTimeout   time.Duration

	p2pCtx context.Context
	swap *swap.Service

	chequeHandler *chequeHandler
}

func New(p2pCtx context.Context, swap *swap.Service, peer swarm.Address, ethAdrr, beneficiary common.Address) (*Bonus, error) {
	wg := &sync.WaitGroup{}

	hbWriteCh := make(chan *writeMsg)
	hbReadCh := make(chan *readMsg)

	addr := Addrs
	dialTimeout := time.Minute * 5
	pingTimeout := time.Second * 10
	retryInterval := time.Second * 10

	activeSession, err := newSession(Dial(addr, retryInterval, dialTimeout))
	if err != nil {
		return nil, err
	}

	hb := newHeartbeater(activeSession, hbWriteCh, hbReadCh)
	ch := newChequeHanler(p2pCtx, peer, swap)


	b := &Bonus{
		peerIDStr:  peer.String(),
		ethAdrrStr: ethAdrr.String(),
		beneficiaryStr: beneficiary.String(),

		peer: peer,
		ethAdrr: ethAdrr,
		beneficiary: beneficiary,

		closeCh: make(chan struct{}),
		wg:      wg,

		hbWriteCh:    hbWriteCh,
		hbReadCh:     hbReadCh,
		activeBeater: hb,
		//ticker: time.NewTicker(time.Second * 15),

		addrs:         Addrs,
		retryInterval: retryInterval,
		dialTimeout:   dialTimeout,
		pingTimeout:   pingTimeout,

		p2pCtx: p2pCtx,
		swap: swap,

		chequeHandler: ch,
	}

	go b.serveHeartbeater()

	return b, nil
}

func (b *Bonus) serveHeartbeater() {
	b.wg.Add(1)
	defer b.wg.Done()

	var hbErr error
	ch := b.chequeHandler

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()
	for {
		select {
		case <-b.closeCh:
			return
		case readMsg := <-b.hbReadCh:
			if readMsg.err != nil {
				fmt.Printf("heartbeater received error: %v\n", readMsg.err)
				hbErr = readMsg.err
				break
			}
			fmt.Printf("heartbeater received CSID %v message: %v\n", readMsg.msg.id, readMsg.msg.msg)

			// todo: add more logic to deal with more kinds of message.
			if err := ch.handleReceivCheque(readMsg.msg.msg); err != nil {
				fmt.Printf("Failed to received cheque to swap service: %v\n", readMsg.msg.msg)
			}
		case t := <-ticker.C:
			fmt.Println("Current time: ", t)
			msg := &message{
				id: CSID_ID_Heartbeat,
				msg: &Heartbeat{
					Peer:    b.peerIDStr,
					EthAddr: b.ethAdrrStr,
					ChequebookAddr: b.beneficiaryStr,
				},
			}

			errCh := make(chan error)

			writeMsg := &writeMsg{
				msg:   msg,
				errCh: errCh,
			}
			b.hbWriteCh <- writeMsg
			err := <-errCh
			if err != nil {
				fmt.Printf("heartbeater write error: %v\n", err)
				hbErr = err
				break
			}
			fmt.Printf("heartbeater write CSID %v message: %v\n", msg.id, msg.msg)
		}

		// close current active heartbeater and rebuild a new one if any error encountered
		if hbErr != nil {
			fmt.Printf("heartbeater err: %v\n", hbErr)
			fmt.Printf("Start newHearbeater.")
			b.activeBeater.close()
			for {
				activeSession, err := newSession(Dial(b.addrs, b.retryInterval, b.dialTimeout))
				if err != nil {
					time.Sleep(time.Second * 5)
					continue
				}

				newHB := newHeartbeater(activeSession, b.hbWriteCh, b.hbReadCh)
				b.activeBeater = newHB
				break
			}
		}
	}
}

// Close close all bonus services in top to down order.
func (b *Bonus) Close() error {
	close(b.closeCh)
	b.wg.Wait()
	return b.activeBeater.close()
}
