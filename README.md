# loadbalancer

Load balancer using ReverseProxy in Golang


Example:
```
func main(){
	serverPool.addNewServer("http://localhost:8080")
	serverPool.addNewServer("http://localhost:8081")
	serverPool.addNewServer("http://localhost:8082")



	http.HandleFunc("/", handleIncomingTraffic)

	fmt.Println("Listening on port 3000.")
	http.ListenAndServe(":3000", nil)
}
```

Currently it sends the requests to all the servers using Round Robin algorithm