package main

import (
	"fmt"
	"gopic/pkg/config"
	"gopic/pkg/server"
)

func main() {
	conf := config.NewConfig("gopic.conf")
	server := server.NewServer(conf)

	fmt.Println("Listening on %s", server.Server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
