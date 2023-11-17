# loadbalancer

Load balancer using ReverseProxy in Golang


Example: refer to `main.go`
```
func main() {
	
	lb := LoadBalancer{}

	lb.serverPool.addNewServer("http://localhost:8080")
	lb.serverPool.addNewServer("http://localhost:8081")
	lb.serverPool.addNewServer("http://localhost:8082")


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//get a server from the serverpool and use it's ServeHTTP method to serve requests
		server := lb.serverPool.getServerRoundRobin()
		server.ReverseProxy.ServeHTTP(w, r)
	})

	fmt.Println("Listening on port 3000.")
	http.ListenAndServe(":3000", nil)
}
```

and then run `go run main.go loadbalancer.go`


Currently it sends the requests to all the servers using Round Robin algorithm