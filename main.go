package main

import (
	"example.com/app/app"
	"example.com/app/database"
)


func init() {
	// create database connection instance for first time
	_ = database.GetInstance()
}

func main() {
	app.Start()
}
