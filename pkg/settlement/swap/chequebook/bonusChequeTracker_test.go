package chequebook

import (
	"github.com/newswarm-lab/new-bee/pkg/logging"
	"github.com/newswarm-lab/new-bee/pkg/statestore/leveldb"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestBonusChequeTracker(t *testing.T) {
	storepath := "./temp/"
	if err := os.RemoveAll(storepath); err != nil {
		t.Fatal(err)
	}

	logger := logging.New(os.Stdout, logrus.TraceLevel)

	store, err := leveldb.NewStateStore(storepath, logger)
	if err != nil {
		t.Fatal(err)
	}

	newTracker := loadBonusChequeTracker(store)

	var chequeKeys = []chequeKeyT{
		"chequeKey0",
		"chequeKey1",
		"chequeKey2",
		"chequeKey3",
		"chequeKey4",
		"chequeKey5",
	}

	for _, key := range chequeKeys {
		newTracker.receiveOneCheque(key)
	}

	if len(newTracker.ChequeKeys) != 6 || newTracker.TotalCheques != 6 || newTracker.CashedIndex != -1 {
		t.Logf("receiveOneCheque performed unexpectedly.")
	}

	cheque, err := newTracker.chequeToCashout()
	if err != nil {
		t.Error(err)
	}
	wanted := chequeKeyT("chequeKey0")
	if cheque != wanted {
		t.Logf("chequeToCashout performed unexpectedly. wanted %q, but goted: %q\n", wanted, cheque)
	}
	if newTracker.CashedIndex != -1 {
		t.Error("CashedIndex shouldn't increase at calling chequeToCashout")
	}

	newTracker.confirmChequeToCashout()
	if newTracker.CashedIndex != 0 {
		t.Error("CashedIndex not increase as expected after calling confirmChequeToCashout")
	}

	if err := newTracker.store(); err != nil {
		t.Error(err)
	}

	loadedTracker := loadBonusChequeTracker(store)
	if len(loadedTracker.ChequeKeys) != loadedTracker.TotalCheques {
		t.Errorf("Loaded bonusChequeTracker unexpectedly.")
	}
	if loadedTracker.CashedIndex != 0 {
		t.Errorf("Loaded bonusChequeTracker unexpectedly.")
	}
	if loadedTracker.TotalCheques != 6 {
		t.Errorf("Loaded bonusChequeTracker unexpectedly.")
	}

	if err := os.RemoveAll(storepath); err != nil {
		t.Fatal(err)
	}
}
