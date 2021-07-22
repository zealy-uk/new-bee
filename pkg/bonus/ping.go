package bonus

import (
	"fmt"
	"github.com/go-ping/ping"
	"sync"
	"time"
)

// Ping pings on given ip addresses and return a map whose addresses are able to receive packets within specified time limit.
func Ping(addrs []string, timeout time.Duration) []string {
	recv := &sync.Map{}
	wg := &sync.WaitGroup{}
	wg.Add(len(addrs))

	addrCh := make(chan string, len(addrs))
	for i := range addrs {
		addrCh <- addrs[i]
		go func() {
			ping_(<-addrCh, recv, timeout)
			wg.Done()
		}()
	}
	wg.Wait()

	var pongAddrs = make([]string, 0, len(addrs))
	recv.Range(func(key, value interface{}) bool {
		pongAddrs = append(pongAddrs, key.(string))
		return true
	})

	return pongAddrs[:len(pongAddrs)]
}

func ping_(addr string, recv *sync.Map, timeout time.Duration) error {
	p, err := ping.NewPinger(addr)
	if err != nil {
		fmt.Printf("Failed to create pinger for :%v\n", addr)
		return err
	}

	p.Timeout = timeout
	p.SetPrivileged(true)

	p.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		p.Stop()
	}

	p.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

		if stats.PacketsRecv >= 1 {
			recv.Store(stats.Addr, struct{}{})
		}
	}

	return p.Run()
}