package main

import (
	"errors"
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

func getDBConf() (*DBConf, error) {

	db := DBConf{}
	var exist bool
	var err error

	if db.host, exist = os.LookupEnv("DB_HOST"); !exist {
		return nil, errors.New("DB_HOST required")
	}
	if portStr, exist := os.LookupEnv("DB_PORT"); !exist {
		return nil, errors.New("DB_PORT required")

	} else {
		db.port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("DB_PORT must be a number %s", err)
		}
	}
	if db.user, exist = os.LookupEnv("DB_USER"); !exist {
		return nil, errors.New("DB_USER required")
	}
	if db.password, exist = os.LookupEnv("DB_PASSWORD"); !exist {
		return nil, errors.New("DB_PASSWORD required")
	}
	if db.name, exist = os.LookupEnv("DB_NAME"); !exist {
		return nil, errors.New("DB_NAME required")
	}

	db.connectionStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		db.host, db.port, db.user, db.password, db.name)

	return &db, nil

}
