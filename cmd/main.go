package main

import (
	"library-backend/config"
	"library-backend/internal/routes"
)

func main() {
	config.InitConfig()

	router := routes.SetupRoutes()

	err := router.Run(":8080")
	if err != nil {
		panic("failed to start server" + err.Error())
	}
}
