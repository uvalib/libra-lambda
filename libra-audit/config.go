package main

import (
	"fmt"
	"os"
	"strconv"
)

// DBConf holds the database cconnection info
type DBConf struct {
	host          string
	port          int
	user          string
	password      string
	name          string
	connectionStr string
}

func getConfig() DBConf {

	db := DBConf{}
	var exist bool

	if db.host, exist = os.LookupEnv("DB_HOST"); !exist {
		panic("DB_HOST required")
	}
	if portStr, exist := os.LookupEnv("DB_PORT"); !exist {
		panic("DB_PORT required")
	} else {
		db.port, _ = strconv.Atoi(portStr)
	}
	if db.user, exist = os.LookupEnv("DB_USER"); !exist {
		panic("DB_USER required")
	}
	if db.password, exist = os.LookupEnv("DB_PASSWORD"); !exist {
		panic("DB_PASSWORD required")
	}
	if db.name, exist = os.LookupEnv("DB_NAME"); !exist {
		panic("DB_NAME required")
	}

	db.connectionStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		db.host, db.port, db.user, db.password, db.name)

	return db

}
