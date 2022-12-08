package main

import (
	"fmt"
	"load-balancer/service"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

func main() {
	domains := []string{"http://localhost:8000", "http://localhost:8001"}
	var pool service.ServicePool
	for _, domain := range domains {
		url, err := url.Parse(domain)
		if err != nil {
			message := fmt.Sprintf("Unable to parse domain: %s", domain)
			panic(message)
		}
		pool.Services = append(pool.Services,
			&service.Service{
				URL:          url,
				Alive:        true,
				ReverseProxy: httputil.NewSingleHostReverseProxy(url)})
	}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		http.HandleFunc("/", pool.Balance)
		fmt.Println("we up")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		log.Println("Starting routine health checks")
		pool.GroupHealthCheck()
	}()
	wg.Wait()
}
