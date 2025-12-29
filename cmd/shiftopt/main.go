package main

import (
	"log"
	"github.com/iannsp/shiftopt/internal/database"
	"github.com/iannsp/shiftopt/internal/scheduler"
)

func main() {
	// 1. Infrastructure
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Data
	database.SeedData(db)

	// 3. Business Logic
	scheduler.RunGreedy(db)
}
