package main

import (
	_ "gorm.io/driver/mysql"
)

func main() {

	server := InitWebServer()
	server.Run(":8080")
}
