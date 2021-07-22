package bonus

import (
"testing"
"time"
)

var (
	addrsForDial = map[string]string{
		"192.168.101.134": "9527",
		"202.182.98.210":  "1398",
		"127.0.0.1":       "44199",
	}

	retryInterval = time.Minute * 5
	dialTimeout   = time.Second * 10
	pingTimeout   = time.Second * 10
)

func TestDial(t *testing.T) {
	conn := Dial(addrsForDial, retryInterval, dialTimeout)
	t.Logf("Connect remote %v\n", conn.RemoteAddr())
}

func TestDialAfterPing(t *testing.T) {
	conn := DialAfterPing(addrsForDial, retryInterval, dialTimeout, pingTimeout)
	t.Logf("Connect remote %v\n", conn.RemoteAddr())
}

func TestAddrsFn(t *testing.T) {url := "testnet.newswarm.info"
	addrPortMp := addrsFn(url, time.Second * 5)
	t.Logf("addrsFn result: %v\n", addrPortMp)
}

func TestWrappedDialAfterPing(t *testing.T)  {
	url := "testapi.newswarm.info"
	addr, conn:= WrappedDialAfterPing(url)
	t.Logf("connect to %v, connection: %v\n", addr, conn)
}