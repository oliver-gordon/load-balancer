package main

import (
	"fmt"
	"load-balancer/service"
	"log"
	"net/url"
	"sync"
)

func main() {
	domains := []string{"http://localhost:8000", "http://localhost:8001", "http://localhost:8002"}
	var pool service.ServicePool
	for _, domain := range domains {
		url, err := url.Parse(domain)
		if err != nil {
			message := fmt.Sprintf("Unable to parse domain: %s", domain)
			panic(message)
		}
		pool.Services = append(pool.Services, &service.Service{URL: url, Alive: false})
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		log.Println("Starting routine health checks")
		pool.GroupHealthCheck()
	}()
	wg.Wait()
}
