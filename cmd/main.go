package main

import "log"

func main() {
	r := router()

	log.Println("Starting listening on port 8181")
	r.Run(":8181")
}