package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Proxy struct {
	localAddr *net.TCPAddr
	balancer  *Balancer
}

func (p *Proxy) Start() {
	listener, err := net.ListenTCP(Network, p.localAddr)
	log.Printf("proxy listening on %v", p.localAddr)

	if err != nil {
		return
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			return
		}
		go p.HandleClient(conn)

	}

}

func (p *Proxy) HandleClient(conn *net.TCPConn) {
	defer conn.Close()
	upstreamTarget := p.balancer.NextAvailableTarget()
	upstreamAddr := upstreamTarget.Address
	upstreamConn, err := net.DialTCP(Network, nil, upstreamAddr)
	if err != nil {
		log.Printf("error connecting to upstream %v", err)
	}
	defer upstreamConn.Close()
	go io.Copy(conn, upstreamConn)
	io.Copy(upstreamConn, conn)
}

func newProxy(port int, balancer *Balancer) (*Proxy, error) {
	localAddr, err := net.ResolveTCPAddr(Network, fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("error resolving downstream address: %v", err)
	}
	return &Proxy{balancer: balancer, localAddr: localAddr}, nil

}

// creates a group of proxy servers pointing to the same loadbalancer
func NewProxyGroup(name string, ports []int, balancer *Balancer) []*Proxy {
	proxyGroups := []*Proxy{}
	for _, port := range ports {
		proxy, err := newProxy(port, balancer)
		if err != nil {
			log.Print(err)
			continue
		}
		go proxy.Start()
		proxyGroups = append(proxyGroups, proxy)
	}
	return proxyGroups
}
