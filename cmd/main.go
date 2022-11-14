package main

import (
	"log"
	"sync"

	"github.com/paularah/proxy/pkg/proxy"
)

func main() {
	var wg sync.WaitGroup
	cfg, err := proxy.LoadConfigFromFile("../config.json")
	if err != nil {
		log.Fatalf("error loading config file %v", err)
	}
	wg.Add(len(cfg.Apps))
	server := proxy.NewServer(cfg)
	server.Bootstrap()
	wg.Wait()

}
