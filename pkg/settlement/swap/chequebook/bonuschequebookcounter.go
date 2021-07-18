package chequebook

import (
	"fmt"
	"github.com/newswarm-lab/new-bee/pkg/storage"
)

const (
	bonusChequebookCounterPrefix = "swap_bonus_chequebook_counter_"
)

type bonusChequebookCounter struct {
	chequebook         chequebookT
	lastReceivedCheque int
	lastCashedCheque   int
}

func initBonusChequebookCounter(chequebook chequebookT, store storage.StateStorer) (*bonusChequebookCounter, error) {
	var b bonusChequebookCounter
	err := store.Get(bonusChequebookCounterKey(chequebook), &b)
	if err != nil {
		if err == storage.ErrNotFound {
			return &bonusChequebookCounter{
				chequebook:         chequebook,
				lastReceivedCheque: -1,
				lastCashedCheque:   -1,
			}, nil
		}
		return nil, err
	}

	return &b, nil
}

func (b *bonusChequebookCounter) receiveOneCheque() *bonusChequebookCounter {
	b.lastReceivedCheque++
	return b
}

// if -1 returned, it implies that temporarily no available cheque for cash out.
func (b *bonusChequebookCounter) chequeToCashout() int {
	if b.lastReceivedCheque < 0 || b.lastCashedCheque == b.lastCashedCheque {
		return -1
	}
	return b.lastCashedCheque + 1
}

func (b *bonusChequebookCounter) confirmChequeToCashout() *bonusChequebookCounter {
	b.lastCashedCheque++
	return b
}

func (b *bonusChequebookCounter) chequesUncashed() int {
	return b.lastReceivedCheque - b.lastCashedCheque
}

func (b *bonusChequebookCounter) store(store storage.StateStorer) error {
	return store.Put(bonusChequebookCounterKey(b.chequebook), b)
}

func bonusChequebookCounterKey(t chequebookT) string {
	return fmt.Sprintf("%s%s", bonusChequebookCounterPrefix, t)
}
