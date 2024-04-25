package main

import (
	"fmt"
	"gotify/internal/server"
)

func main() {
	server := server.NewServer()

	fmt.Println("Server is running!", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		panic("cannot start server")
	}

}
