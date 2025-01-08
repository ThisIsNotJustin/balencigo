package loadbalancer

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type HTTPServer struct {
	addr   string
	proxy  http.Handler
	status bool
	mu     sync.Mutex
}

type Server interface {
	Address() string
	IsActive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

func CreateServerUtil(addr string, ctx context.Context) *HTTPServer {
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

/*
func CreateServerReverseProxy(addr string, ctx context.Context, lb *LoadBalancer) *HTTPServer {
	serverURL, err := url.Parse(addr)
	if err != nil {
		log.Printf("Invalid Server Address %q: %v", addr, err)
		return nil
	}

	server := &HTTPServer{
		addr:   addr,
		proxy:  CreateReverseProxy(lb, serverURL),
		status: false,
	}

	server.StartHealthCheck(10*time.Second, ctx)
	return server
}
*/

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
