package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type HTTPServer struct {
	addr   string
	proxy  *httputil.ReverseProxy
	status bool
	mu     sync.Mutex
}

type Server interface {
	Address() string
	IsActive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

type LoadBalancer struct {
	port            string
	roundRobinIndex int
	servers         []Server
	mu              sync.Mutex
}

func CreateLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinIndex: 0,
		servers:         servers,
	}
}

func CreateServer(addr string, ctx context.Context) *HTTPServer {
	serverURL, err := url.Parse(addr)
	if err != nil {
		log.Printf("Invalid Server Address %q: %v", addr, err)
		return nil
	}

	server := &HTTPServer{
		addr:   addr,
		proxy:  httputil.NewSingleHostReverseProxy(serverURL),
		status: false,
	}

	server.StartHealthCheck(10*time.Second, ctx)
	return server
}

func (server *HTTPServer) Address() string {
	return server.addr
}

func (server *HTTPServer) IsActive() bool {
	server.mu.Lock()
	defer server.mu.Unlock()

	return server.status
}

func (server *HTTPServer) Serve(rw http.ResponseWriter, req *http.Request) {
	server.proxy.ServeHTTP(rw, req)
}

func (server *HTTPServer) StartHealthCheck(interval time.Duration, ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Stopping Health Checks for server %q", server.addr)
				return

			case <-time.After(interval):
				status := server.CheckHealth()
				server.mu.Lock()
				server.status = status
				server.mu.Unlock()
			}
		}
	}()
}

func (server *HTTPServer) CheckHealth() bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	for i := 0; i < 3; i++ {
		resp, err := client.Get(server.addr)
		if err == nil {
			resp.Body.Close()
			return resp.StatusCode >= 200 && resp.StatusCode < 300
		}
		time.Sleep(1000 * time.Millisecond)
	}

	log.Printf("Server, %q, failed\n", server.addr)
	return false
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

	return nil, fmt.Errorf("No available servers\n")
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

func InitializeServers(addresses []string, ctx context.Context) []Server {
	var servers []Server
	for _, addr := range addresses {
		server := CreateServer(addr, ctx)
		if server != nil {
			servers = append(servers, server)
		}
	}

	return servers
}

func main() {
	addresses := []string{
		"https://github.com",
		"https://duckduckgo.com",
		"https://kagi.com",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	servers := InitializeServers(addresses, ctx)
	if len(servers) == 0 {
		log.Fatalf("No Server Addresses Available")
	}

	lb := CreateLoadBalancer("8080", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.ServeProxy(rw, req)
	}

	http.HandleFunc("/", handleRedirect)
	log.Printf("Serving Requests at localhost:%s \n", lb.port)

	err := http.ListenAndServe(":"+lb.port, nil)
	if err != nil {
		log.Fatalf("Critical error: %v", err)
	}
}
