package main

import "moneytransfer-api/api"

func main() {
	server := api.NewServer()
	server.Start(9000)
}
