package main

import (
	"context"
	"fmt"
	"time"
	"wordofwisdom/internal/server"
	"wordofwisdom/internal/storage"
)

func main() {
	serverConfig := server.ServerConfig{
		Port: 9000,
		Host: "0.0.0.0",
	}

	storage := storage.InitInMemoryStorage(time.Now())

	ctx := context.Background()
	ctx = context.WithValue(ctx, "storage", storage)
	ctx = context.WithValue(ctx, "storageExpiration", int64(3600))

	err := server.Run(ctx, serverConfig)
	if err != nil {
		fmt.Println(err)
	}

}
