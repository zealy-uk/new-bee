package chequebook

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/newswarm-lab/new-bee/pkg/swarm"
	"math/big"
)

const (
	peerPrefix                = "swap_chequebook_peer_" // copied from pkg/settlement/swap/addressbook.go
	uncashedBonusChequePrefix = "swap_bonus_chequebook_uncashed_cheque_"
	cashedBonusChequePrefix   = "swap_bonus_chequebook_cashed_cheque_"
)

type chequebookT string
type chequeKeyT string
type cashedChequeKeyT string
type chequeTxHashT string

func bonusReceivedChequeKey(cheque *SignedCheque) chequeKeyT {
	return chequeKeyT(fmt.Sprintf("%schequebook:%x_chequeid:%s_amout:%d", uncashedBonusChequePrefix, cheque.Chequebook, cheque.Id, cheque.Beneficiary))
}
func bonusCashedChequeKey(cheque *SignedCheque, txHash common.Hash) cashedChequeKeyT {
	return cashedChequeKeyT(fmt.Sprintf("%schequebook:%x_chequeid:%s_amout:%d_txHash:%s", cashedBonusChequePrefix, cheque.Chequebook, cheque.Id, cheque.Beneficiary, txHash))
}

type BonousChequeStore struct {
	keyCache        map[chequebookT][]chequeKeyT
	keyTxCache      map[chequebookT]map[cashedChequeKeyT]chequeTxHashT
	txKeyCache      map[chequebookT]map[chequeTxHashT]cashedChequeKeyT
	chequebookCache []chequebookT

	*ChequeStoreImp
}

// NewChequeStore creates new BonousChequeStore
func NewBonusChequeStore(cs *ChequeStoreImp) *BonousChequeStore {

	var chequebookCache []chequebookT

	if err := cs.store.Iterate(peerPrefix, func(_, value []byte) (stop bool, err error) {
		var chequebook common.Address

		if err := json.Unmarshal(value, &chequebook); err != nil {
			return true, fmt.Errorf("invalid chequebook value %q: %w", string(value), err)
		}

		chequebookCache = append(chequebookCache, chequebookT(chequebook.String()))
		return false, nil
	}); err != nil {
		panic(fmt.Errorf("iteration failed: build chequebook cache from storage: %w\n", err))
	}


}

// earliestBonusCheque returns the earliest received but not cashed signed cheque and its corresponding key.
func (r *BonousChequeStore) EarliestBonusCheque(chequebook chequebookT) (chequeKeyT, *SignedCheque, error) {
	keysqueue, ok := r.keyCache[chequebook]
	if !ok || len(keysqueue) < 1 {
		return "", nil, ErrNoCheque
	}

	chequekey := keysqueue[0]

	var cheque *SignedCheque

	if err := r.store.Get(string(chequekey), cheque); err != nil {
		return "", nil, err
	}

	return chequekey, cheque, nil
}

// DeleteEearliestBonusCheque deletes the earliest received signed cheque from both storage and cache.
// This function should called only after the signed cheque has been cashed out successfully.
func (r *BonousChequeStore) DeleteEearliestBonusCheque(chequebook chequebookT) error {
	keysqueue, ok := r.keyCache[chequebook]
	if ok {
		if len(keysqueue) > 0 {
			chequekey := keysqueue[0]

			if err := r.store.Delete(string(chequekey)); err != nil {
				return nil
			}

			keysqueue = keysqueue[1:]
		}
	}

	return nil
}

// ReceiveBonusCheque stores given signed cheque and caches its key.
func (r *BonousChequeStore) StoreReceivedBonusCheque(cheque *SignedCheque) (*big.Int, error) {
	// verify we are the beneficiary
	if cheque.Beneficiary != r.beneficiary {
		return nil, ErrWrongBeneficiary
	}

	// don't allow concurrent processing of cheques
	// this would be sufficient on a per chequebookT basis
	r.lock.Lock()
	defer r.lock.Unlock()

	receivedChequeKey := bonusReceivedChequeKey(cheque)
	if err := r.store.Put(string(receivedChequeKey), cheque); err != nil {
		return nil, err
	}

	chequebook := chequebookT(cheque.Chequebook.String())
	keyscache, ok := r.keyCache[chequebook]
	if !ok {
		keyscache = make([]chequeKeyT, 0, 256)
		r.chequebookCache = append(r.chequebookCache, chequebook)
	}

	keyscache = append(keyscache, receivedChequeKey)

	fmt.Printf("%v received and cached\n", receivedChequeKey)
	return cheque.CumulativePayout, nil
}

// StoreCashedBonusCheque stores given already cashed signed cheque and caches key-txhash/txhas-key map.
func (r *BonousChequeStore) StoreCashedBonusCheque(cheque *SignedCheque, txhash common.Hash) error {
	cashedChequeKey_ := bonusCashedChequeKey(cheque)
	if err := r.store.Put(string(cashedChequeKey_), &cashoutAction{
		TxHash: txhash,
		Cheque: *cheque,
	}); err != nil {
		return err
	}

	chequebook := chequebookT(cheque.Chequebook.String())
	keyTxMp, ok := r.keyTxCache[chequebook]
	if !ok {
		keyTxMp = make(map[cashedChequeKeyT]chequeTxHashT, 256)
	}

	txHash := chequeTxHashT(txhash.String())
	keyTxMp[cashedChequeKey_] = txHash

	txKeyMp, ok := r.txKeyCache[chequebook]
	if !ok {
		txKeyMp = make(map[chequeTxHashT]cashedChequeKeyT, 256)
	}

	txKeyMp[txHash] = cashedChequeKey_

	return nil
}

// ReceivedBonusCheques returns all received bonus cheques to a chequebookT that not yet cashed.
func (r *BonousChequeStore) ReceivedBonusCheques(chequebook chequebookT) ([]*SignedCheque, error) {
	keysQueue, ok := r.keyCache[chequebook]
	if !ok {
		return nil, nil
	}

	var results []*SignedCheque
	for _, key := range keysQueue {
		var cheque *SignedCheque
		if err := r.store.Get(string(key), cheque); err != nil {
			return results, err
		}

		results = append(results, cheque)
	}

	return results, nil
}

// ReceivedBonusCheques returns all received bonus cheques that not yet cashed.
func (r *BonousChequeStore) AllReceivedBonusCheques() (map[chequebookT][]*SignedCheque, error) {
	var results = make(map[chequebookT][]*SignedCheque, 256)
	for _, chequebook := range r.chequebookCache {
		cheques, err := r.ReceivedBonusCheques(chequebook)

		results[chequebook] = cheques
		if err != nil {
			return results, nil
		}
	}
	return results, nil
}
