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

	// 3. Compare Strategies

    // Strategy A: illegal but cheap 
	scheduler.RunGreedy(db)

	// Strategy B: The Real World (Legal but expensive)
	scheduler.RunConstrained(db)

	// Strategy B+safety. The Safety Constraint (Max 8h + 1 Senior)
	scheduler.RunSafeSchedule(db)
}
