package chequebook

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	//peerPrefix                = "swap_chequebook_peer_" // copied from pkg/settlement/swap/addressbook.go
	uncashedBonusChequePrefix = "swap_bonus_chequebook_uncashed_cheque_"
	cashedBonusChequePrefix   = "swap_bonus_chequebook_cashed_cheque_"
)

type chequebookT string
type chequeKeyT string
type cashedChequeKeyT string
//type chequeTxHashT string

func bonusReceivedChequeKey(chequebook chequebookT, chequeId int64) chequeKeyT {
	return chequeKeyT(fmt.Sprintf("%schequebook:%s_chequeid:%d", uncashedBonusChequePrefix, chequebook, chequeId))
}
func bonusCashedChequeKey(chequebook chequebookT, chequeId int64) cashedChequeKeyT {
	return cashedChequeKeyT(fmt.Sprintf("%schequebook:%s_chequeid:%d", cashedBonusChequePrefix, chequebook, chequeId))
}

//func bonusReceivedChequeKeyPrefix(chequebook chequebookT) string {
//	return fmt.Sprintf("%schequebook:%s", uncashedBonusChequePrefix, chequebook)
//}
//
//func bonusCashedChequeKeyPrefix(chequebook chequebookT) cashedChequeKeyT {
//	return cashedChequeKeyT(fmt.Sprintf("%schequebook:%s", cashedBonusChequePrefix, chequebook))
//}

type BonousChequeStore struct {
	//chequebooksCache  []chequebookT
	//uncashedKeysCache map[chequebookT][]chequeKeyT
	//cashedKeysCache   map[chequebookT][]cashedChequeKeyT

	//keyTxCache      map[chequebookT]map[cashedChequeKeyT]chequeTxHashT
	//txKeyCache      map[chequebookT]map[chequeTxHashT]cashedChequeKeyT

	*ChequeStoreImp

	chequebookcounters map[chequebookT]*bonusChequebookCounter
}

// NewBonusChequeStore creates new BonousChequeStore
func NewBonusChequeStore(cs *ChequeStoreImp) *BonousChequeStore {
	//
	//var chequebookCache = make([]chequebookT, 0, 32)
	//
	//if err := cs.store.Iterate(peerPrefix, func(_, value []byte) (stop bool, err error) {
	//	var chequebook common.Address
	//
	//	if err := json.Unmarshal(value, &chequebook); err != nil {
	//		return true, fmt.Errorf("invalid chequebook value %q: %w", string(value), err)
	//	}
	//
	//	chequebookCache = append(chequebookCache, chequebookT(chequebook.String()))
	//	return false, nil
	//}); err != nil {
	//	panic(fmt.Errorf("iteration failed: build chequebook cache from storage: %w\n", err))
	//}
	//
	//var keysCache = make(map[chequebookT][]chequeKeyT, 32)
	//
	//for _, chequebook := range chequebookCache {
	//	keysCache[chequebook] = make([]chequeKeyT, 0, 128)
	//	prefix := bonusReceivedChequeKeyPrefix(chequebook)
	//	if err := cs.store.Iterate(prefix, func(key, _ []byte) (stop bool, err error) {
	//		keysCache[chequebook] = append(keysCache[chequebook], chequeKeyT(string(key)))
	//		return false, nil
	//	}); err != nil {
	//		panic(fmt.Errorf("iteration failed: build cheque keys cache from storage: %w\n", err))
	//	}
	//}

	return &BonousChequeStore{cs, make(map[chequebookT]*bonusChequebookCounter, 64)}

}

// ChequeToCashout returns the earliest received but not cashed signed cheque
func (r *BonousChequeStore) ChequeToCashout(chequebook chequebookT) (*SignedCheque, error) {
	chequebookCounter, err := r.chequebookCounter(chequebook)
	if err != nil {
		return nil, err
	}

	chequeId, err := chequebookCounter.chequeToCashout()
	if err != nil {
		return nil, err
	}

	chequeKey := bonusReceivedChequeKey(chequebook, chequeId)

	var cheque SignedCheque

	if err := r.store.Get(string(chequeKey), &cheque); err != nil {
		return nil, err
	}

	return &cheque, nil
}

//// DeleteEearliestBonusCheque deletes the earliest received signed cheque from both storage and cache.
//// This function should called only after the signed cheque has been cashed out successfully.
//func (r *BonousChequeStore) DeleteEearliestBonusCheque(chequebook chequebookT) error {
//	keysqueue, ok := r.uncashedKeysCache[chequebook]
//	if ok {
//		if len(keysqueue) > 0 {
//			chequekey := keysqueue[0]
//
//			if err := r.store.Delete(string(chequekey)); err != nil {
//				return nil
//			}
//
//			keysqueue = keysqueue[1:]
//		}
//	}
//
//	return nil
//}

// StoreReceivedBonusCheque stores given signed cheque and caches its key.
func (r *BonousChequeStore) StoreReceivedBonusCheque(cheque *SignedCheque) (*big.Int, error) {
	// verify we are the beneficiary
	//if cheque.Beneficiary != r.beneficiary {
	//	return nil, ErrWrongBeneficiary
	//}

	// don't allow concurrent processing of cheques
	// this would be sufficient on a per chequebookT basis
	r.lock.Lock()
	defer r.lock.Unlock()

	chequebook := chequebookT(cheque.Chequebook.String())
	chequebookCounter, err := r.chequebookCounter(chequebook)
	if err != nil {
		return nil, err
	}


	receivedChequeKey := bonusReceivedChequeKey(chequebook, cheque.Id.Int64())
	if err := r.store.Put(string(receivedChequeKey), cheque); err != nil {
		return nil, err
	}

	if err := chequebookCounter.receiveOneCheque().store(r.store); err != nil {
		return nil, err
	}

	fmt.Printf("%v received and stored.\n", receivedChequeKey)
	return cheque.CumulativePayout, nil
}

// StoreCashedBonusCheque stores given already cashed signed cheque and caches key-txhash/txhas-key map.
func (r *BonousChequeStore) StoreCashedBonusCheque(cheque *SignedCheque, txhash common.Hash) error {
	chequebook := chequebookT(cheque.Chequebook.String())
	chequeId := cheque.Id.Int64()
	cashedChequeKey_ := bonusCashedChequeKey(chequebook, chequeId)
	if err := r.store.Put(string(cashedChequeKey_), &cashoutAction{
		TxHash: txhash,
		Cheque: *cheque,
	}); err != nil {
		return err
	}

	chequebookCounter, err := r.chequebookCounter(chequebook)
	if err != nil {
		return err
	}

	if err := chequebookCounter.confirmChequeToCashout().store(r.store); err != nil {
		return err
	}


	return nil
}

//// ReceivedBonusCheques returns all received bonus cheques to a chequebookT that not yet cashed.
//func (r *BonousChequeStore) ReceivedBonusCheques(chequebook chequebookT) ([]*SignedCheque, error) {
//	keysQueue, ok := r.uncashedKeysCache[chequebook]
//	if !ok {
//		return nil, nil
//	}
//
//	var results []*SignedCheque
//	for _, key := range keysQueue {
//		var cheque SignedCheque
//		if err := r.store.Get(string(key), &cheque); err != nil {
//			return results, err
//		}
//
//		results = append(results, cheque)
//	}
//
//	return results, nil
//}
//
//// AllReceivedBonusCheques returns all received bonus cheques that not yet cashed.
//func (r *BonousChequeStore) AllReceivedBonusCheques() (map[chequebookT][]*SignedCheque, error) {
//	var results = make(map[chequebookT][]*SignedCheque, 256)
//	for _, chequebook := range r.chequebooksCache {
//		cheques, err := r.ReceivedBonusCheques(chequebook)
//
//		results[chequebook] = cheques
//		if err != nil {
//			return results, nil
//		}
//	}
//	return results, nil
//}

func (r *BonousChequeStore) chequebookCounter(chequebook chequebookT) (*bonusChequebookCounter, error)  {
	_, ok := r.chequebookcounters[chequebook]
	if !ok {
		counter, err := initBonusChequebookCounter(chequebook, r.store);
		if err != nil {
			return nil, err
		}
		r.chequebookcounters[chequebook] = counter
	}

	return r.chequebookcounters[chequebook], nil
}