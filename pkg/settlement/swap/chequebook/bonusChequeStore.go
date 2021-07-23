package chequebook

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/newswarm-lab/new-bee/pkg/bonus/log"
	"github.com/newswarm-lab/new-bee/pkg/storage"
	"math/big"
	"sync"
)

const (
	uncashedBonusChequePrefix = "bonus_uncashed_cheque_"
	cashedBonusChequePrefix   = "bonus_cashed_cheque_"
)

type chequebookT string
type chequeKeyT string
type cashedChequeKeyT string

func bonusReceivedChequeKey(chequebook chequebookT, chequeId int64) chequeKeyT {
	return chequeKeyT(fmt.Sprintf("%schequebook:%s_chequeid:%d", uncashedBonusChequePrefix, chequebook, chequeId))
}
func bonusCashedChequeKey(chequebook chequebookT, chequeId int64) cashedChequeKeyT {
	return cashedChequeKeyT(fmt.Sprintf("%schequebook:%s_chequeid:%d", cashedBonusChequePrefix, chequebook, chequeId))
}

type BonousChequeStore struct {
	tracker *bonusChequeTracker
	lock    *sync.Mutex
	storer  storage.StateStorer
}

// NewBonusChequeStore creates new BonousChequeStore
func newBonusChequeStore(lock *sync.Mutex, storer storage.StateStorer) *BonousChequeStore {
	return &BonousChequeStore{
		tracker: loadBonusChequeTracker(storer),
		lock:    lock,
		storer:  storer,
	}
}

// ChequeToCashout returns the earliest received but not cashed signed cheque
func (r *BonousChequeStore) ChequeToCashout() (*SignedCheque, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	chequeK, err := r.tracker.chequeToCashout()
	if err != nil {
		log.Error("chequeToCashout() failed. Err: %v\n", err)
		return nil, err
	}

	var cheque SignedCheque

	if err := r.storer.Get(string(chequeK), &cheque); err != nil {
		log.Error("BonusStateStorer.Get() failed. Err:%v\n", err)
		return nil, err
	}

	return &cheque, nil
}

// StoreReceivedBonusCheque stores given signed cheque and caches its key.
func (r *BonousChequeStore) StoreReceivedBonusCheque(cheque *SignedCheque) (*big.Int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	chequebook := chequebookT(cheque.Chequebook.String())
	chequeKey := bonusReceivedChequeKey(chequebook, cheque.Id.Int64())

	if err := r.storer.Put(string(chequeKey), cheque); err != nil {
		log.Error("failed to store cheque:%q\n", chequeKey)
		return nil, err
	}

	if err := r.tracker.receiveOneCheque(chequeKey).store(); err != nil {
		log.Error("failed to store bonusChequeTracker.\n")
		return nil, err
	}

	return cheque.CumulativePayout, nil
}

// StoreCashedBonusCheque stores given already cashed signed cheque and caches key-txhash/txhas-key map.
func (r *BonousChequeStore) StoreCashedBonusCheque(cheque *SignedCheque, txhash common.Hash) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	chequebook := chequebookT(cheque.Chequebook.String())
	chequeId := cheque.Id.Int64()
	cashedChequeKey_ := bonusCashedChequeKey(chequebook, chequeId)
	if err := r.storer.Put(string(cashedChequeKey_), &cashoutAction{
		TxHash: txhash,
		Cheque: *cheque,
	}); err != nil {
		log.Error("failed to store cashed bonus cheque %q.\n", cashedChequeKey_)
		return err
	}

	if err := r.tracker.confirmChequeToCashout().store(); err != nil {
		log.Error("failed to store bonusChequeTracker.\n")
		return err
	}

	return nil
}

func (r *BonousChequeStore) BonusReceivedUncashedCheques() ([]*SignedCheque, error)  {
	keys := r.tracker.uncashedChequeKey()
	var cheques []*SignedCheque
	for _, key := range keys {
		var cheque SignedCheque
		if err := r.storer.Get(string(key), &cheque); err != nil {
			if err == storage.ErrNotFound {
				continue
			}
			return nil, err
		}
		cheques = append(cheques, &cheque)
	}
	return cheques, nil
}