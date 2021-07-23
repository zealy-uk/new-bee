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

			return &btracker_

		}
		panic(fmt.Errorf("failed to load bonusChequeTracker from storage. Err: %w\n", err))
	}

	return &bonusChequeTracker{
		tracker: &tracker_,
		storer:  storer,
	}
}

func (b *bonusChequeTracker) receiveOneCheque(chequeK chequeKeyT) *bonusChequeTracker {
	b.ChequeKeys = append(b.ChequeKeys, chequeK)
	b.TotalCheques++
	return b
}

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

func (b *bonusChequeTracker) uncashedChequeKey() []chequeKeyT {
	if !b.uncashedChequeExists()  {
		return nil
	}
	return b.ChequeKeys[b.CashedIndex+1:]
}

func (b *bonusChequeTracker) uncashedChequeExists() bool {
	return b.TotalCheques > 0 && b.CashedIndex < b.TotalCheques-1
}