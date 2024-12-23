# BalenciGo

BalenciGo is a lightweight, scalable load balancer written in Go.

## Features
- Round-Robin Load Balancing: Distributes incoming requests evenly across available servers.
- Health Checks: Coninuously monitors server health to avoid routing to inactive or failing servers.
- Concurrency: Built keeping synchronization for thread-safe operations in mind by utilizing a mutex.

## Getting Started
### Prerequisites
- Go
- Internet Connection (for testing functionality)

## Installation
Clone the repository:
```bash
git clone https://github.com/ThisIsNotJustin/balencigo.git
```
Run the project:
```bash
go run main.go
```

## Example Usage
By default BalenciGo routes requests to the following demo servers
- GitHub
- DuckDuckGo
- Kagi

## Testing
To test the load balancer, visit http://localhost:8080 in the browser. 
Requests will then be forwarded to the configured servers in a round-robin cycle.

## Configuration
Edit the addresses array in main.go to specify the server URLs:
```go
addresses := []string {
    "https://github.com",
	"https://duckduckgo.com",
	"https://kagi.com",
}
```

Edit the time interval of the health checks by editing the following
```go
server.StartHealthCheck(10*time.Second, ctx)
```