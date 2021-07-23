package network

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

func dial() (string, net.Conn, error) {
	addrs, err := getAddrs("http://testapi.newswarm.info:10080/nodelist")
	if err != nil {
		return "", nil, err
	}

	type addrT struct {
		addr string
		port string
	}
	addrsD := make([]addrT, 0, len(addrs))
	for addr, port := range addrs {
		addrsD = append(addrsD, addrT{addr, port})
	}

	rand.Seed(86)
	for len(addrsD) > 0 {
		i := rand.Intn(len(addrsD))
		addr := addrsD[i].addr + ":" + addrsD[i].port

		conn, err := net.DialTimeout("tcp", addr, time.Second*10)
		if err != nil {
			addrsD = append(addrsD[0:i], addrsD[i+1:]...)
			continue
		}
		return addr, conn, nil
	}

	return "", nil, errors.New("failed to connect any endpoint")
}

func getAddrs(url string) (map[string]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	addrsByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var addrsList []string
	if err := json.Unmarshal(addrsByte, &addrsList); err != nil {
		return nil, err
	}

	if len(addrsList) < 1 {
		return nil, errors.New("no address response")
	}

	var resMp = make(map[string]string, 64)
	for _, addrS := range addrsList {
		addrPort := strings.Split(addrS, ":")
		resMp[addrPort[0]] = addrPort[1]
	}

	return resMp, nil
}
