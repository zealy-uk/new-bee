package bonus

import (
	"context"
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/newswarm-lab/new-bee/pkg/settlement/swap"
	"github.com/newswarm-lab/new-bee/pkg/settlement/swap/chequebook"
	"github.com/newswarm-lab/new-bee/pkg/swarm"
)

type chequeHandler struct {
	p2pCtx context.Context
	peer swarm.Address
	swap *swap.Service
}

func newChequeHanler(p2pCtx context.Context, peer swarm.Address, swap *swap.Service) *chequeHandler {
	return &chequeHandler{
		p2pCtx: p2pCtx,
		peer: peer,
		swap: swap,
	}
}

func (c *chequeHandler) handleReceivCheque(msg proto.Message) error  {
	ec := msg.(*EmitCheque)

	signedCheck := &chequebook.SignedCheque{}
	if err := json.Unmarshal(ec.Cheque, signedCheck); err != nil {
		return err
	}

	return c.swap.ReceiveBonusCheque(c.p2pCtx, c.peer, signedCheck)
}
