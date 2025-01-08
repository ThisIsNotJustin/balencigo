package loadbalancer

import (
	"net/http"
	"net/url"
)

type ReverseProxy struct {
	lb  *LoadBalancer
	url *url.URL
}

func CreateReverseProxy(lb *LoadBalancer, url *url.URL) *ReverseProxy {
	return &ReverseProxy{
		lb:  lb,
		url: url,
	}
}

func (rp *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	/*
		log.Printf("Reverse Proxy receiving requests\n")

		if rp.url != nil {
			log.Printf("Forwarding requests to %q", rp.url.String())
		}

		target, err := rp.lb.NextAvailableServer()
		if err != nil {
			log.Printf("Error getting next available server: %v", err)
			http.Error(rw, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}

		URL, err := url.Parse(target.Address())
		if err != nil {
			log.Printf("Error parsing address: %v", err)
			http.Error(rw, "Bad Gateway", http.StatusBadGateway)
			return
		}
	*/

}

/*

/*

	Simple Reverse Proxy

		rp := CreateReverseProxy(lb)
		log.Printf("Reverse Proxy receiving requests\n")
		target, err := lb.NextAvailableServer()
		if err != nil {
			log.Fatalf("Error getting next available server")
		}

		URL, err := url.Parse(target.Address())
		if err != nil {
			log.Fatalf("Bad Gateway")
		}
		proxy := httputil.NewSingleHostReverseProxy(URL)
		proxy.ServeHTTP(rw, req)
		log.Printf("Serving Requests\n")

		err := http.ListenAndServe(":8081", reverseProxy)
		if err != nil {
			log.Fatalf("Critical error: %v", err)
		}
*/
