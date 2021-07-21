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
)

type bonusChequeTracker struct {
	*tracker

	storer storage.StateStorer
}

type tracker struct {
	TotalCheques int
	CashedIndex  int
	ChequeKeys   []chequeKeyT
}

func loadBonusChequeTracker(storer storage.StateStorer) *bonusChequeTracker {
	var tracker_ tracker
	err := storer.Get(bonusChequeTrackerKey, &tracker_)
	if err != nil {
		if err == storage.ErrNotFound {
			btracker_ := bonusChequeTracker{
				tracker: &tracker{
					TotalCheques: 0,
					CashedIndex:  -1,
					ChequeKeys:   make([]chequeKeyT, 0, 1024),
				},

				storer: storer,
			}

			fmt.Printf("✅✅✅✅✅ new bonusChequeTracker: %+#v\n", btracker_)
			return &btracker_

		}
		panic(fmt.Errorf("xxxxxxxxxx failed to load bonusChequeTracker from storage. Err: %w\n", err))
	}

	fmt.Printf("✅✅✅✅✅ loaded bonusChequeTracker. TotalCheques: %v, CashedIndex: %v\n", tracker_.TotalCheques, tracker_.CashedIndex)
	return &bonusChequeTracker{
		tracker: &tracker_,
		storer:  storer,
	}
}

func (b *bonusChequeTracker) receiveOneCheque(chequeK chequeKeyT) *bonusChequeTracker {
	b.ChequeKeys = append(b.ChequeKeys, chequeK)
	b.TotalCheques++
	fmt.Printf("✅✅✅✅✅ cheque %q cached. totalCheques=%d, cashedIndex=%d\n", chequeK, b.TotalCheques, b.CashedIndex)
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

func (b *bonusChequeTracker) store() error {
	if err := b.storer.Put(bonusChequeTrackerKey, b); err != nil {
		return err
	}
	return nil
}
