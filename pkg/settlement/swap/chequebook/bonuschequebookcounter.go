package chequebook

import (
	"errors"
	"fmt"
	"github.com/newswarm-lab/new-bee/pkg/storage"
	"sync"
)

const (
	bonusChequebookCounterPrefix = "swap_bonus_chequebook_counter_"
)

var (
	ErrNoCashableCheque = errors.New("no cheque can be cashed out")


	once = &sync.Once{}
	defaultBonusChequebookCounter *bonusChequebookCounter
)

type bonusChequebookCounter struct {
	chequebook chequebookT
	chequeKeys []chequeKeyT
}


func initBonusChequebookCounter(chequebook chequebookT, store storage.StateStorer) *bonusChequebookCounter {
	once.Do(func() {
		fmt.Printf("started once.Do()")
		var chequeKeys []chequeKeyT
		err := store.Get(bonusChequebookCounterKey(chequebook), chequeKeys)
		if err != nil {
			if err == storage.ErrNotFound {
				fmt.Printf("coudn't find chequeKeys for chequebook %q,\n and a new bonusChequebookCounter will be created.\n", chequebook)

				d := &bonusChequebookCounter{
					chequebook: chequebook,
					chequeKeys: make([]chequeKeyT, 0, 1024),
				}

				defaultBonusChequebookCounter = d
			}
			fmt.Printf("failed to load bonusChequebookCounter from storage. Err: %w\n", err)
			panic(fmt.Errorf("failed to load bonusChequebookCounter from storage. Err: %w\n", err))
		}
	})

	fmt.Printf("current defaultBonusChequebookCounter status: chequebook: %q, chequeKeys length: %v\n", defaultBonusChequebookCounter.chequebook, len(defaultBonusChequebookCounter.chequeKeys))
	return defaultBonusChequebookCounter
}

func (b *bonusChequebookCounter) receiveOneCheque(chequeK chequeKeyT) *bonusChequebookCounter {
	//b.lastReceivedCheque++

	b.chequeKeys = append(b.chequeKeys, chequeK)
	fmt.Printf("bonuschequebookcounter cash. chequebook: %v, chequekey:%v\n", b.chequebook, chequeK)
	fmt.Printf("bonusChequebookCounter: chequeBook: %v; chequeKeys length:%v\n", b.chequebook, len(b.chequeKeys))
	return b
}

// if "" returned, it implies that temporarily no available cheque for cash out.
func (b *bonusChequebookCounter) chequeToCashout() (chequeKeyT, error) {
	if len(b.chequeKeys) < 1 {
		return "", ErrNoCashableCheque
	}

	chequeK := b.chequeKeys[0]
	return chequeK, nil
}

// todo: improve logic
func (b *bonusChequebookCounter) confirmChequeToCashout() *bonusChequebookCounter {
	b.chequeKeys = b.chequeKeys[1:]
	return b
}

func (b *bonusChequebookCounter) store(store storage.StateStorer) error {
	if err := store.Put(bonusChequebookCounterKey(b.chequebook), b.chequeKeys); err != nil {
		fmt.Printf("failed to store bonuschequebookcounter: %v. ERROR: %v\n", b.chequebook, err)
		return err
	}
	fmt.Printf("bonuschequebookcounter: %v stored\n", b.chequebook)
	return nil
}

func bonusChequebookCounterKey(t chequebookT) string {
	return fmt.Sprintf("%s%s", bonusChequebookCounterPrefix, t)
}
