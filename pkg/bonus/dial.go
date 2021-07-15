package bonus

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// todo: configure addresses of interest
var Addrs = map[string]string{
	"192.168.101.134": "9527",
	"127.0.0.1":       "44199",
}

// Dial randomly picks an address from given addresses for dialing.
// Dial returns a connection immediately as long as connected successfully.
// Dial will repeat the process generally depicted as above until one connection established.
func Dial(addrs map[string]string, retryInterval, dialTimeout time.Duration) net.Conn {
	type addrT struct {
		addr string
		port string
	}
	addrsD := make([]addrT, 0, len(addrs))
	for addr, port := range addrs {
		addrsD = append(addrsD, addrT{addr, port})
	}

	for {
		rand.Seed(86)
		for len(addrsD) > 0 {
			i := rand.Intn(len(addrsD))
			addr := addrsD[i].addr + ":" + addrsD[i].port

			fmt.Printf("start dailing %v\n", addr)
			conn, err := net.DialTimeout("tcp", addr, dialTimeout)
			if err != nil {
				fmt.Printf("failed to dial %v. ERROR: %v\n", addr, err)
				addrsD = append(addrsD[0:i], addrsD[i+1:]...)
				continue
			}

			fmt.Printf("successfully to connect %v\n", addr)
			return conn
		}

		fmt.Printf("failed to dial any endpoint. \n Program will retry after %v\n", retryInterval)
		<-time.After(retryInterval)
	}
}

// DialAfterPing performs like Dial but it works on addresses that Ping successfully.
func DialAfterPing(addrs map[string]string, retryInterval, dialTimeout, pingTimeout time.Duration) net.Conn {
	addrsP := make([]string, 0, len(addrs))
	for addr, _ := range addrs {
		addrsP = append(addrsP, addr)
	}

	for {
		fmt.Println("start Ping process")
		pongAddrs := Ping(addrsP, pingTimeout)

		if len(pongAddrs) < 1 {
			fmt.Printf("ping with no response and dial process won't start.\n Program will retry after %v\n", retryInterval)
			<-time.After(retryInterval)
			continue
		}

		rand.Seed(86)
		for len(pongAddrs) > 0 {
			i := rand.Intn(len(pongAddrs))
			addr := pongAddrs[i]

			addr = addr + ":" + addrs[addr]
			conn, err := net.DialTimeout("tcp", addr, dialTimeout)
			if err != nil {
				fmt.Printf("failed to dial %v. ERROR: %v\n", addr, err)
				pongAddrs = append(pongAddrs[0:i], pongAddrs[i+1:]...)
				continue
			}

			fmt.Printf("successfully to connect %v\n", addr)
			return conn
		}

		fmt.Printf("failed to dial any endpoint even after ping success. \n Program will retry after %v\n", retryInterval)
		<-time.After(retryInterval)
	}
}
