package main

import (
	"github.com/vds/RestaurantManagement/pkg/database/mysql"
	"github.com/vds/RestaurantManagement/pkg/server"
)

func main(){
	// create database instance
	db, err := mysql.NewMySqlDB()
	if err != nil {
		panic(err)
	}

	// create server
	s, err := server.NewServer(db)
	if err != nil {
		panic(err)
	}
	if err := s.Start(":5000"); err != nil {
		panic(err)
	}
}