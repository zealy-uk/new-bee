package network

import (
	"fmt"
	"testing"
)

const testUrl = "http://testapi.newswarm.info:10080/nodelist"

func TestGetAddrs(t *testing.T) {
	addrsMp, err := getAddrs(testUrl)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(addrsMp)
}

func TestDial(t *testing.T) {
	addr, conn, err := dial()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("connected to: %q\nconnection: %v\n", addr, conn)
}
