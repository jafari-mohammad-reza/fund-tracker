package main

import (
	"github.com/jafari-mohammad-reza/fund-tracker/api"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	data.SetupRedisClient()
	api.NewServer()
}
