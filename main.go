package main

import (
	"fmt"
	"net/http"
)


func main() {
	// lb := LoadBalancer{balancerType: "round-robin"}
	lb := LoadBalancer{balancerType: "least-connections"}

	lb.addNewServer("http://localhost:8080")
	lb.addNewServer("http://localhost:8081")
	lb.addNewServer("http://localhost:8082")

	http.HandleFunc("/", lb.handleIncomingRequests)

	fmt.Println("Listening on port 3000.")
	http.ListenAndServe(":3000", nil)
}
