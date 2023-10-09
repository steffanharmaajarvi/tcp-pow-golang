package main

import (
	"context"
	"fmt"
	"wordofwisdom/internal/client"
	"wordofwisdom/internal/server"
)

func main() {
	fmt.Println("start client")

	serverConfig := server.ServerConfig{
		Port: 9000,
		Host: "server",
	}

	// init context to pass config down
	ctx := context.Background()

	address := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)

	// run client
	err := client.Run(ctx, address)
	if err != nil {
		fmt.Println("client error:", err)
	}

}
