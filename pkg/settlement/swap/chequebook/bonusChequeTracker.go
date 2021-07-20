package chequebook

import (
	"errors"
	"fmt"
	"github.com/newswarm-lab/new-bee/pkg/storage"
)

const (
	bonusChequeTrackerKey = "swap_bonus_cheque_tracker"
)

var (
	ErrNoCashableCheque = errors.New("no cheque can be cashed out")

	//myBonusChequeTracker *bonusChequeTracker
)

type bonusChequeTracker struct {
	TotalCheques int
	CashedIndex  int
	ChequeKeys   []chequeKeyT
}

func initbonusChequeTracker(store storage.StateStorer) *bonusChequeTracker {
		var tracker bonusChequeTracker
		err := store.Get(bonusChequeTrackerKey, &tracker)
		if err != nil {
			if err == storage.ErrNotFound {
				tracker = bonusChequeTracker{
					ChequeKeys:   make([]chequeKeyT, 0, 1024),
					TotalCheques: 0,
					CashedIndex:  -1,
				}

				fmt.Printf("new bonusChequeTracker: %+#v\n", tracker)
				return &tracker

			}
			panic(fmt.Errorf("failed to load bonusChequeTracker from storage. Err: %w\n", err))
		}

	fmt.Printf("loaded bonusChequeTracker: %+v\n", tracker)
	return &tracker
}

func (b *bonusChequeTracker) receiveOneCheque(chequeK chequeKeyT) *bonusChequeTracker {
	b.ChequeKeys = append(b.ChequeKeys, chequeK)
	b.TotalCheques++
	fmt.Printf("cheque %q cached. total cheques=%d, cashed index=%d\n", chequeK, b.TotalCheques, b.CashedIndex)
	return b
}

// if "" returned, it implies that temporarily no available cheque for cash out.
func (b *bonusChequeTracker) chequeToCashout() (chequeKeyT, error) {
	if b.TotalCheques < 1 || b.CashedIndex == b.TotalCheques-1 {
		return "", ErrNoCashableCheque
	}

	chequeK := b.ChequeKeys[b.CashedIndex+1]
	return chequeK, nil
}

func (b *bonusChequeTracker) confirmChequeToCashout() *bonusChequeTracker {
	b.CashedIndex++
	return b
}

func (b *bonusChequeTracker) store(store storage.StateStorer) error {
	if err := store.Put(bonusChequeTrackerKey, b); err != nil {
		fmt.Printf("failed to store: %+v. ERROR: %v\n", b, err)
		return err
	}
	fmt.Printf("stored %+v\n", b)
	return nil
}
