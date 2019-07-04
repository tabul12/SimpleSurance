package main

import (
	"fmt"
	cserver "simpleSurance/Server"
)

func main() {
	server := &cserver.Server{Addr : ":8080", Pattern: "/view/"}
	err := server.Start()
	if err != nil {
		fmt.Println("Servers could not start with reason: %s", err)
	}
}