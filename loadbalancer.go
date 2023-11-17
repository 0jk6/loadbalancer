package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

//store the backend server and it's reverse proxy object
type Server struct {
	url string
	ReverseProxy *httputil.ReverseProxy
}

//store the servers here in the list
type ServerPool struct {
	servers []*Server //store servers in a list
	mutex sync.Mutex //for concurrent reads and writes from server pool
}


//global variables
var serverPool ServerPool
var roundRobinIndex int //by default this is zero


//method to create a reverseproxy for a given server
func (sp *ServerPool) createReverseProxy(serverURL string) *httputil.ReverseProxy {
	//extract origin from serverURL
	origin, _ := url.Parse(serverURL)
	
	//create a director for httputil.ReverseProxy struct
	director := func(r *http.Request){
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", origin.Host)
		r.URL.Scheme = "http"
		r.URL.Host = origin.Host
	}

	reverseProxy := &httputil.ReverseProxy{Director: director}

	return reverseProxy
}

func (sp *ServerPool) addNewServer(serverURL string) {
	server := &Server{
		url: serverURL,
		ReverseProxy: sp.createReverseProxy(serverURL),
	}

	sp.servers = append(sp.servers, server)
}

func (sp *ServerPool) getServerRoundRobin() *Server {
	//lock the mutex before reading
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	server := sp.servers[roundRobinIndex%len(sp.servers)]
	roundRobinIndex++;

	return server
}


//following function will handle incoming traffic to the load balancer
func handleIncomingTraffic(w http.ResponseWriter, r *http.Request){
	//get a server from the serverpool
	server := serverPool.getServerRoundRobin()
	//use the ServeHTTP method attached to the reverseproxy to server requests
	server.ReverseProxy.ServeHTTP(w, r)
}

func main(){
	serverPool.addNewServer("http://localhost:8080")
	serverPool.addNewServer("http://localhost:8081")
	serverPool.addNewServer("http://localhost:8082")



	http.HandleFunc("/", handleIncomingTraffic)

	fmt.Println("Listening on port 3000.")
	http.ListenAndServe(":3000", nil)
}