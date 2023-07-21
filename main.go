package main

import (
	"github.com/jafari-mohammad-reza/fund-tracker/api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	api.NewServer()

}
