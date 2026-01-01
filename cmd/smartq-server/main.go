package main

import (
	"github.com/smartq/smartq/internal/api"
)

func main() {
	router := api.NewRouter()
	router.Run(":8080")
}

