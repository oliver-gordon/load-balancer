package service

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type ServicePool struct {
	Services []*Service
	current  uint32
}

type Service struct {
	URL          *url.URL
	Alive        bool
	failedCalls  uint32
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (sp *ServicePool) Balance(w http.ResponseWriter, r *http.Request) {
	as := sp.GetService()
	fmt.Println(as)
	if as != nil {
		as.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Unavailable", http.StatusServiceUnavailable)
}

func (sp *ServicePool) GetService() *Service {
	nextIndex := int(atomic.AddUint32(&sp.current, uint32(1)) % uint32(len(sp.Services)))
	service := sp.Services[nextIndex]
	if service.serviceIsAvailable() {
		sp.setCurrent(nextIndex)
		return service
	}
	log.Println("No service available")
	return nil
}

func (sp *ServicePool) setCurrent(index int) {
	atomic.StoreUint32(&sp.current, uint32(index))
}

// Mark if failedCalls limit is reached, mark the host domain as void.
func (s *Service) increaseFailedCalls() {
	atomic.AddUint32(&s.failedCalls, 1)
}
