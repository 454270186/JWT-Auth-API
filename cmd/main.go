package main

import (
	"jwtAuth/domain"
	"log"
)

func main() {
	db, err := domain.NewDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := router(db)

	log.Println("Starting listening on port 8181")
	r.Run(":8181")
}