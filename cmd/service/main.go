package main

import (
	"fmt"
	"motus/goechoservice/internal/server"
	"net/http"
)

func main() {
	fmt.Println("[INFO] Starting Thrift service on port 9090")
	err := http.ListenAndServe(":9090", server.New())
	if err != nil {
		panic(err)
	}
}
