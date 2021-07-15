package bonus

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/proto"
	"github.com/newswarm-lab/new-bee/pkg/settlement/swap/chequebook"
	"github.com/newswarm-lab/new-bee/pkg/storage"
)

const (
	// prefix for the persistence key
	ChequePrefix = "cheque"
)

var hEntry = handleEntry{
	CSID_ID_EmitCheque: handleReceiveCheque,
}

var stateStore storage.StateStorer

type hanndleFn func(msg proto.Message) error
type handleEntry map[CSID]hanndleFn

// handleEmitCheque stores cheque book using underlying storage.StateStorer.
func handleReceiveCheque(msg proto.Message) error {
	ec := msg.(*EmitCheque)

	sc := &chequebook.SignedCheque{}
	if err := json.Unmarshal(ec.Cheque, sc); err != nil {
		return err
	}

	// store the accepted cheque
	err := stateStore.Put(chequeKey(sc.Chequebook, sc.ID), sc)
	if err != nil {
		return err
	}
	return nil
}


// chequeKey computes the key where to store cheque received from a chequebook.
func chequeKey(chequebook common.Address, chequeID int) string {
	return fmt.Sprintf("%s_%x_%d", ChequePrefix, chequebook, chequeID)
}

func initStateStorer(storer storage.StateStorer) {
	stateStore = storer
}
