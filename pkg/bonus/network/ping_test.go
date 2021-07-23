package network

import (
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	addrs := []string{
		"www.google.com",
		"cn.bing.com",
		"202.182.98.210",

		"192.168.101.134",
		"127.0.0.1",
	}

	pongAddrs := Ping(addrs, time.Second*10)

	t.Logf("pong addresses: %v\n", pongAddrs)
}