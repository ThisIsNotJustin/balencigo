package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ThisIsNotJustin/balencigo/loadbalancer"
)

func InitializeServers(addresses []string, ctx context.Context) []loadbalancer.Server {
	var servers []loadbalancer.Server
	for _, addr := range addresses {
		server := loadbalancer.CreateServerUtil(addr, ctx)
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

	lb := loadbalancer.CreateLoadBalancer("8080", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.ServeProxy(rw, req)
	}

	http.HandleFunc("/", handleRedirect)
	log.Printf("Serving Requests at localhost:%s \n", lb.GetPort())

	err := http.ListenAndServe(":"+lb.GetPort(), nil)
	if err != nil {
		log.Fatalf("Critical error: %v", err)
	}
}
