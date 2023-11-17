package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

//store the backend server and it's reverse proxy object
type Server struct {
	url          string
	ReverseProxy *httputil.ReverseProxy
}

//store the servers here in the list
type ServerPool struct {
	servers []*Server  //store servers in a list
	mutex   sync.Mutex //for concurrent reads and writes from server pool
}

//global variables
var roundRobinIndex int //by default this is zero

//method to create a reverseproxy for a given server
func (sp *ServerPool) createReverseProxy(serverURL string) *httputil.ReverseProxy {
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

func (sp *ServerPool) addNewServer(serverURL string) {
	server := &Server{
		url:          serverURL,
		ReverseProxy: sp.createReverseProxy(serverURL),
	}

	sp.servers = append(sp.servers, server)
}

//returns a server from the server pool in a round robin fashion
func (sp *ServerPool) getServerRoundRobin() *Server {
	//lock the mutex before reading
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	server := sp.servers[roundRobinIndex%len(sp.servers)]
	roundRobinIndex++

	return server
}
