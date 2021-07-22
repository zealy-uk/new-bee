package bonus

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

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

			log.Printf("start dailing %v\n", addr)
			conn, err := net.DialTimeout("tcp", addr, dialTimeout)
			if err != nil {
				log.Printf("failed to dial %v. ERROR: %v\n", addr, err)
				addrsD = append(addrsD[0:i], addrsD[i+1:]...)
				continue
			}

			log.Printf("successfully to connect %v\n", addr)
			return conn
		}

		log.Printf("failed to dial any endpoint. \n Program will retry after %v\n", retryInterval)
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
		log.Println("start Ping process")
		pongAddrs := Ping(addrsP, pingTimeout)

		if len(pongAddrs) < 1 {
			log.Printf("ping with no response and dial process won't start.\n Program will retry after %v\n", retryInterval)
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
				log.Printf("failed to dial %v. ERROR: %v\n", addr, err)
				pongAddrs = append(pongAddrs[0:i], pongAddrs[i+1:]...)
				continue
			}

			log.Printf("successfully to connect %v\n", addr)
			return conn
		}

		log.Printf("failed to dial any endpoint even after ping success. \n Program will retry after %v\n", retryInterval)
		<-time.After(retryInterval)
	}
}

type addrsFnT func(url string, retryInterval time.Duration) map[string]string

// DialAfterPingWithDynamicAddrs performs like DialAfterPing but it depends on addrsFn to retry refreshed addresses.
func DialAfterPingWithDynamicAddrs(url string, addrsFn addrsFnT, retryInterval, dialTimeout, pingTimeout time.Duration) (string, net.Conn) {
retrieveAddrs:
	addrs := addrsFn(url, retryInterval)
	addrsP := make([]string, 0, len(addrs))
	for addr, _ := range addrs {
		addrsP = append(addrsP, addr)
	}

	//for {
		log.Println("start Ping process")
		pongAddrs := Ping(addrsP, pingTimeout)

		if len(pongAddrs) < 1 {
			log.Printf("ping with no response and dial process won't start.\n Program will retrieveAddrs after %v\n", retryInterval)
			<-time.After(retryInterval)
			goto retrieveAddrs
		}

		rand.Seed(86)
		for len(pongAddrs) > 0 {
			i := rand.Intn(len(pongAddrs))
			addr := pongAddrs[i]

			addr = addr + ":" + addrs[addr]
			conn, err := net.DialTimeout("tcp", addr, dialTimeout)
			if err != nil {
				log.Printf("failed to dial %v. ERROR: %v\n", addr, err)
				pongAddrs = append(pongAddrs[0:i], pongAddrs[i+1:]...)
				continue
			}

			log.Printf("successfully to connect %v\n", addr)
			return addr, conn
		}

		log.Printf("failed to dial any endpoint even after ping success. \n Program will retrieveAddrs after %v\n", retryInterval)
		<-time.After(retryInterval)
		goto retrieveAddrs
	//}
}

// WrappedDialAfterPing wraps DialAfterPingWithDynamicAddrs function
func WrappedDialAfterPing(url string) (string, net.Conn) {
	return DialAfterPingWithDynamicAddrs(url, addrsFn, time.Minute*5, time.Second*10, time.Second*10)
}

func addrsFn(url string, retryInterval time.Duration) map[string]string {
	var addrsList []string
	for {
		res, err := http.Get(url)
		if err != nil {
			log.Printf("failed to get addresses list. Error: %v\n", err)
			<-time.After(retryInterval)
			continue
		}

		addrsByte, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("failed to read addresse bytes from respond body. Error: %v\n", err)
		}

		if err := json.Unmarshal(addrsByte, &addrsList); err != nil {
			log.Fatalf("failed to unmarshal responded addresses bytes. Error: %v\n", err)
		}

		_ = res.Body.Close()
		break
	}

	var res = make(map[string]string, 64)
	for _, addrS := range addrsList {
		addrPort := strings.Split(addrS, ":")
		res[addrPort[0]] = addrPort[1]
	}

	return res
}


