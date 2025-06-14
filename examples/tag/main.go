package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anvidev/goenv"
)

type serverConfig struct {
	Env      string `goenv:"ENV"`
	Api      apiConfig
	Database databaseConfig
}

type apiConfig struct {
	Port         int           `goenv:"API_PORT,default=8080"`
	ReadTimeout  time.Duration `goenv:"API_READ_TIMEOUT"`
	WriteTimeout time.Duration `goenv:"API_WRITE_TIMEOUT"`
}

type databaseConfig struct {
	Name       string `goenv:"DB_NAME,required"`
	ConnString string `goenv:"DB_CONN_STRING,required"`
}

func main() {
	var config serverConfig

	if err := goenv.Struct(&config); err != nil {
		log.Fatal(err)
	}

	// Config is now populated
	fmt.Println(config.Env)
	fmt.Println(config.Api.Port)
}
