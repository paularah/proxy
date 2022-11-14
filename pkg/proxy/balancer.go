package proxy

import (
	"net"
	"sync/atomic"
	"time"
)

const (
	DialTimeout = 10
	Network     = "tcp"
)

type UpstreamTarget struct {
	Address *net.TCPAddr
}

type Balancer struct {
	upstreamTargets []UpstreamTarget
	count           uint64
}

// checks that a target backend works
func (t *UpstreamTarget) IsAlive() bool {
	conn, err := net.DialTimeout(Network, t.Address.String(), (time.Second * DialTimeout))
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

func newUpstreamTarget(name string) (UpstreamTarget, error) {
	var upstreamTarget UpstreamTarget
	address, err := net.ResolveTCPAddr(Network, name)
	if err != nil {
		return upstreamTarget, err
	}
	upstreamTarget.Address = address
	return upstreamTarget, err
}

func NewBalancer(targets []string) *Balancer {
	upstreamTargetList := []UpstreamTarget{}
	for _, target := range targets {
		upstreamTarget, err := newUpstreamTarget(target)
		if err != nil {
			continue
		}
		if ok := upstreamTarget.IsAlive(); ok != true {
			continue
		}

		upstreamTargetList = append(upstreamTargetList, upstreamTarget)
	}
	return &Balancer{upstreamTargets: upstreamTargetList}
}

func (b *Balancer) Next() int {
	return int(atomic.AddUint64(&b.count, uint64(1)))
}

func (b *Balancer) NextAvailableTarget() *UpstreamTarget {
	for tries := 0; tries < len(b.upstreamTargets); tries++ {
		next := b.Next() % len(b.upstreamTargets)
		if ok := b.upstreamTargets[next].IsAlive(); ok {
			return &b.upstreamTargets[next]
		}
	}
	return nil
}
