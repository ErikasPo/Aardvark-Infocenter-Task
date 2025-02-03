package main

import (
	"Infocenter/Application"
	"Infocenter/Infrastructure"
)

func main() {
	service := application.NewMessageService()
	infrastructure.StartServer(service)
}
