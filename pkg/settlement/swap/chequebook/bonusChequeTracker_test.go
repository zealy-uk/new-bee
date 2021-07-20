package chequebook

import (
	"github.com/newswarm-lab/new-bee/pkg/logging"
	"github.com/newswarm-lab/new-bee/pkg/statestore/leveldb"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestBonusChequeTracker(t *testing.T) {
	logger := logging.New(os.Stdout, logrus.TraceLevel)

	store, err := leveldb.NewStateStore("./temp/", logger)
	if err != nil {
		t.Error(err)
	}

	_ = store
}