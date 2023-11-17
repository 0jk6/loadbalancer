package main

import (
	// "fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

//global variables
var roundRobinIndex int //by default this is zero

//store the backend server and it's reverse proxy object
type Server struct {
	url          string
	ReverseProxy *httputil.ReverseProxy
	activeConnections int
	mutex sync.Mutex
}

//store the servers here in the list
type ServerPool struct {
	servers []*Server  //store servers in a list
	mutex   sync.Mutex //for concurrent reads and writes from server pool
}

//create a load balancer struct to wrap the ServerPool
//we can later add type of load balancing in here
type LoadBalancer struct {
	serverPool ServerPool
	balancerType string
}

//handle incoming requests
func (lb *LoadBalancer) handleIncomingRequests(w http.ResponseWriter, r *http.Request){
	server := lb.getServer()

	server.mutex.Lock()
	server.activeConnections++
	server.mutex.Unlock()

	server.ReverseProxy.ServeHTTP(w, r)

	server.mutex.Lock()
	server.activeConnections--
	server.mutex.Unlock()
}

//method to create a reverseproxy for a given server
func (lb *LoadBalancer) createReverseProxy(serverURL string) *httputil.ReverseProxy {
	//extract origin from serverURL
	origin, _ := url.Parse(serverURL)

	//create a director for httputil.ReverseProxy struct
	director := func(r *http.Request) {
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", origin.Host)
		r.URL.Scheme = "http"
		r.URL.Host = origin.Host
	}

	reverseProxy := &httputil.ReverseProxy{Director: director}

	return reverseProxy
}

func (lb *LoadBalancer) addNewServer(serverURL string) {
	server := &Server{
		url:          serverURL,
		ReverseProxy: lb.createReverseProxy(serverURL),
		activeConnections: 0,
	}

	lb.serverPool.servers = append(lb.serverPool.servers, server)
}

//returns a server from the server pool in a round robin fashion
func (sp *ServerPool) getServerRoundRobin() *Server {
	//lock the mutex before reading
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	server := sp.servers[roundRobinIndex%len(sp.servers)]
	roundRobinIndex++

	// fmt.Printf("serving from %s\nactive connections: %d\n", server.url, server.activeConnections)

	return server
}


//returns the server with least connections
func (sp *ServerPool) getServerLeastConnections() *Server {{
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	leastConnectionsServer := sp.servers[0].activeConnections
	serverIndex := 0

	for i, server := range sp.servers {

		//lock and unlock mutex while accessing the activeConnections variable
		server.mutex.Lock()
		if server.activeConnections < leastConnectionsServer {
			leastConnectionsServer = server.activeConnections
			serverIndex = i
		}
		server.mutex.Unlock()
	}

	server := sp.servers[serverIndex]

	// fmt.Printf("serving from %s\nactive connections: %d\n", server.url, server.activeConnections)

	return server
}}


func (lb *LoadBalancer) getServer() *Server {
	if lb.balancerType == "round-robin" {
		return lb.serverPool.getServerRoundRobin()
	} else if lb.balancerType == "least-connections" {
		return lb.serverPool.getServerLeastConnections()
	}

	panic("balancer not found")
}