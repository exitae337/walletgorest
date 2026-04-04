package main

import "github.com/exitae337/walletgorest/internal/config"

func main() {
	// 1. Read and Init config file
	cfg := config.MustLoad()
	// 2. Init logger for REST app
	// 3. Init database connection
	// 4. Init Server and Handlers
	// 5. Start server
}
