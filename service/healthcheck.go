package service

import (
	"log"
	"net"
	"time"
)

func (sp *ServicePool) GroupHealthCheck() {
	var status string
	for {
		for _, service := range sp.Services {
			alive := service.isServiceUp()
			if !alive {
				status = "\u001b[31mdown\u001b[0m"
				service.setServiceDown()
				service.increaseFailedCalls()
			} else {
				status = "\u001b[32mup\u001b[0m"
				service.setServiceUp()
			}
			log.Printf("status: %s - [%s]", service.URL, status)
		}
		// TODO: make sleep duration dependent on an an input variable
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

func (s *Service) setServiceDown() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Alive = false
}

func (s *Service) setServiceUp() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Alive = true
}

func (s *Service) serviceIsAvailable() bool {
	s.mux.RLock()
	available := s.Alive
	s.mux.RUnlock()
	return available
}
