package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/8legd/hugocms/server"
)

// example server implementation using env. vars. for configuration
func main() {
	acc := os.Getenv("HUGOCMS_ACC")
	pwd := os.Getenv("HUGOCMS_PWD")
	port, err := strconv.Atoi(os.Getenv("HUGOCMS_PRT"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		//TODO more graceful exit!
	}

	server.ListenAndServe(
		port,
		server.Auth{acc, pwd},
		server.DB_MySQL, // for sqlite use server.DB_SQLite
	)
}
