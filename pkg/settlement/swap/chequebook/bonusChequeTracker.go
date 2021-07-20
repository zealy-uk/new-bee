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

	myBonusChequeTracker *bonusChequeTracker
)

type bonusChequeTracker struct {
	ChequeKeys   []chequeKeyT
	TotalCheques int
	CashedIndex  int
}

func initbonusChequeTracker(chequebook chequebookT, store storage.StateStorer) *bonusChequeTracker {
	if myBonusChequeTracker == nil {
		//fmt.Printf("initbonusChequeTracker")
		//var tracker bonusChequeTracker
		//err := store.Get(bonusChequeTrackerKey, &tracker)
		//if err != nil {
		//	if err == storage.ErrNotFound {
		//		fmt.Printf("coudn't find ChequeKeys for Chequebook %q, a new bonusChequeTracker will be created.\n", chequebook)
		//
				myBonusChequeTracker = &bonusChequeTracker{
					ChequeKeys:   make([]chequeKeyT, 0, 1024),
					TotalCheques: 0,
					CashedIndex:  -1,
				}

		//		fmt.Printf("Init new bonus cheque tracker: %+#v\n", myBonusChequeTracker)
		//
		//	}
		//	fmt.Printf("failed to load bonusChequeTracker from storage. Err: %v\n", err)
		//	panic(fmt.Errorf("failed to load bonusChequeTracker from storage. Err: %w\n", err))
		//}
	}

	fmt.Printf("current myBonusChequeTracker: %+v\n", myBonusChequeTracker)
	return myBonusChequeTracker
}

func (b *bonusChequeTracker) receiveOneCheque(chequeK chequeKeyT) *bonusChequeTracker {
	b.ChequeKeys = append(b.ChequeKeys, chequeK)
	b.TotalCheques++
	fmt.Printf("bonusChequeTracker cash. Chequebook: %+#v, chequekey:%v\n", b, chequeK)
	fmt.Printf("bonusChequeTracker: chequeBook: %+#v; ChequeKeys length:%v\n", b, len(b.ChequeKeys))
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
		fmt.Printf("failed to store bonusChequeTracker: %+#v. ERROR: %v\n", b, err)
		return err
	}
	fmt.Printf("bonusChequeTracker: %#+v stored\n", b)
	return nil
}
