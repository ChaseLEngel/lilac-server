package main

import (
	"fmt"
	"github.com/chaselengel/lilac/config"
	"github.com/chaselengel/lilac/jwt"
	"github.com/chaselengel/lilac/logger"
	"github.com/gorilla/handlers"
	"net/http"
)

var log *logger.Logger
var jwtData *jwt.JwtData
var conf config.Config

func main() {
	var err error
	conf, err = config.Parse("./config.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse config: %v", err))
	}
	log = logger.New(conf.LogFile)

	if initDatabase(conf.Database) != nil {
		panic("Failed to init database")
	}

	groups, err := allGroups()
	if err != nil {
		panic(fmt.Sprintf("Failed to init checker with groups.", err))
	}
	InitChecker(groups)
	defer master.Stop()

	jwtData, err = jwt.Init(conf.PrivateKey, conf.PublicKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to init jw: %v", err))
	}

	router := NewRouter()
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"Authorization"})
	credentials := handlers.AllowCredentials()
	validators := handlers.AllowedOriginValidator(
		func(origin string) bool {
			return true
		},
	)
	panic(http.ListenAndServe(conf.Port, handlers.CORS(methods, origins, headers, validators, credentials)(router)))
}
