package service

import (
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

type ServicePool struct {
	Services []*Service
}

type Service struct {
	URL          *url.URL
	Alive        bool
	howManyPings int8
	mux          sync.RWMutex
}

func (sp *ServicePool) GroupHealthCheck() {
	var status string
	for {
		for _, service := range sp.Services {

			alive := service.isServiceUp()
			if !alive {
				status = "down"
				service.SetServiceDown()
			} else {
				status = "up"
				service.SetServiceUp()
			}

			log.Printf("status: %s - [%s]", service.URL, status)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *Service) isServiceUp() bool {
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", s.URL.Host, timeout)
	if err != nil {
		return false
	}
	err = conn.Close()
	if err != nil {
		log.Fatal("Unable to close TCP connection, error: ", err)
	}
	return true
}

func (s *Service) SetServiceDown() {
	s.mux.Lock()
	s.Alive = true
	s.mux.Unlock()

}
func (s *Service) SetServiceUp() {
	s.mux.Lock()
	s.Alive = false
	s.mux.Unlock()
}
