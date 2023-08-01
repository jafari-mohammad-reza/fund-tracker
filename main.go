package main

import (
	"github.com/jafari-mohammad-reza/fund-tracker/api"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/joho/godotenv"
)

// @title Fund Tracker API
// @version 1.0
// @description API for retrieving fund market data.
// @termsOfService http://swagger.io/terms/
// @contact.name Mohammadreza jafari
// @contact.url http://www.swagger.io/support
// @contact.email mohammadrezajafari.dev@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:5000
// @BasePath /api/v1
// @schemes http
func main() {
	godotenv.Load(".env")
	data.SetupRedisClient()
	api.NewServer()
}
