package proxy

import (
	"fmt"
	"syscall"
)

type ProxyGroup struct {
	name    string
	proxies []*Proxy
}

type Server struct {
	config Config
	apps   []ProxyGroup
}

func (s *Server) Bootstrap() {
	setUlimit()
	for _, app := range s.config.Apps {
		balancer := NewBalancer(app.Targets)
		proxies := NewProxyGroup(app.Name, app.Ports, balancer)
		s.apps = append(s.apps, ProxyGroup{name: app.Name, proxies: proxies})
	}
}

func setUlimit() error {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		return err
	}
	rLimit.Cur = rLimit.Max
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}

func (s *Server) Start() {
	for _, app := range s.apps {
		for _, proxy := range app.proxies {
			fmt.Println(app.name, proxy.localAddr.String())
			go proxy.Start()
		}
	}
}

func NewServer(config Config) *Server {
	return &Server{config: config}

}
