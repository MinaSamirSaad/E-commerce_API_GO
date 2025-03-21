package main

import (
	"log"

	"github.com/MinaSamirSaad/ecommerce/cmd/api"
	"github.com/MinaSamirSaad/ecommerce/config"
	"github.com/MinaSamirSaad/ecommerce/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	// connect to db
	dataBase, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal("error in connecting to db")
	}
	db.InitStorage(dataBase)

	server := api.NewAPIServer(config.Envs.Port, dataBase)
	if err := server.Run(); err != nil {
		log.Fatal("error in running the server")
	}
}
