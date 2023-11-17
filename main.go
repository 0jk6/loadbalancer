package main

import (
	"fmt"
	"net/http"
)


func main() {
	var serverPool ServerPool

	serverPool.addNewServer("http://localhost:8080")
	serverPool.addNewServer("http://localhost:8081")
	serverPool.addNewServer("http://localhost:8082")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//get a server from the serverpool and use it's ServeHTTP method to serve requests
		server := serverPool.getServerRoundRobin()
		server.ReverseProxy.ServeHTTP(w, r)
	})

	fmt.Println("Listening on port 3000.")
	http.ListenAndServe(":3000", nil)
}
