package chequebook

import (
	"errors"
	"fmt"
	"github.com/newswarm-lab/new-bee/pkg/storage"
)

const (
	bonusChequeKeyCachePrefix = "swap_bonus_chequebook_counter_"
)

var (
	ErrNoCashableCheque = errors.New("no cheque can be cashed out")

	defaultbonusChequeKeyCache *bonusChequeKeyCache
)

type bonusChequeKeyCache struct {
	chequebook chequebookT
	chequeKeys []chequeKeyT
}


func initBonusChequeKeyCache(chequebook chequebookT, store storage.StateStorer) *bonusChequeKeyCache {
	if defaultbonusChequeKeyCache == nil {
		fmt.Printf("initBonusChequeKeyCache")
		var chequeKeys []chequeKeyT
		err := store.Get(bonusChequeKeyCacheKey(chequebook), chequeKeys)
		if err != nil {
			if err == storage.ErrNotFound {
				fmt.Printf("coudn't find chequeKeys for chequebook %q,\n and a new bonusChequeKeyCache will be created.\n", chequebook)

				d := &bonusChequeKeyCache{
					chequebook: chequebook,
					chequeKeys: make([]chequeKeyT, 0, 1024),
				}

				defaultbonusChequeKeyCache = d
			}
			fmt.Printf("failed to load bonusChequeKeyCache from storage. Err: %w\n", err)
			panic(fmt.Errorf("failed to load bonusChequeKeyCache from storage. Err: %w\n", err))
	}

	}

	fmt.Printf("current defaultbonusChequeKeyCache status: chequebook: %q, chequeKeys length: %v\n", defaultbonusChequeKeyCache.chequebook, len(defaultbonusChequeKeyCache.chequeKeys))
	return defaultbonusChequeKeyCache
}

func (b *bonusChequeKeyCache) receiveOneCheque(chequeK chequeKeyT) *bonusChequeKeyCache {
	//b.lastReceivedCheque++

	b.chequeKeys = append(b.chequeKeys, chequeK)
	fmt.Printf("bonusChequeKeyCache cash. chequebook: %v, chequekey:%v\n", b.chequebook, chequeK)
	fmt.Printf("bonusChequeKeyCache: chequeBook: %v; chequeKeys length:%v\n", b.chequebook, len(b.chequeKeys))
	return b
}

// if "" returned, it implies that temporarily no available cheque for cash out.
func (b *bonusChequeKeyCache) chequeToCashout() (chequeKeyT, error) {
	if len(b.chequeKeys) < 1 {
		return "", ErrNoCashableCheque
	}

	chequeK := b.chequeKeys[0]
	return chequeK, nil
}

// todo: improve logic
func (b *bonusChequeKeyCache) confirmChequeToCashout() *bonusChequeKeyCache {
	b.chequeKeys = b.chequeKeys[1:]
	return b
}

func (b *bonusChequeKeyCache) store(store storage.StateStorer) error {
	if err := store.Put(bonusChequeKeyCacheKey(b.chequebook), b.chequeKeys); err != nil {
		fmt.Printf("failed to store bonusChequeKeyCache: %v. ERROR: %v\n", b.chequebook, err)
		return err
	}
	fmt.Printf("bonusChequeKeyCache: %v stored\n", b.chequebook)
	return nil
}

func bonusChequeKeyCacheKey(t chequebookT) string {
	return fmt.Sprintf("%s%s", bonusChequeKeyCachePrefix, t)
}
