package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type LoadBalancer struct {
	port            string
	roundRobinIndex int
	servers         []Server
	mu              sync.Mutex
}

func (lb *LoadBalancer) GetPort() string {
	return lb.port
}

func CreateLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinIndex: 0,
		servers:         servers,
	}
}

func (lb *LoadBalancer) NextAvailableServer() (Server, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i := 0; i < len(lb.servers); i++ {
		server := lb.servers[lb.roundRobinIndex%len(lb.servers)]
		lb.roundRobinIndex++

		if server.IsActive() {
			return server, nil
		}

		log.Printf("Server, %q, is in active", server.Address())
	}

	return nil, fmt.Errorf("no available servers")
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	target, err := lb.NextAvailableServer()
	if err != nil {
		http.Error(rw, "Service Unavailable", http.StatusServiceUnavailable)
		log.Printf("No available servers to handle the request")
		return
	}

	log.Printf("Forwarding request to %q\n", target.Address())
	target.Serve(rw, req)
}
