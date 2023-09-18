package main

import "github.com/YoussefMahmod/MoneyTransfering-API/api"

func main() {
	server := api.NewServer()
	server.Start(9000)
}
