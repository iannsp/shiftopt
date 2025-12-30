package main

import (
	"log"

	"github.com/iannsp/shiftopt/internal/database"
	"github.com/iannsp/shiftopt/internal/scheduler"
)

func main() {
	db, err := database.InitDB()
	if err != nil { log.Fatal(err) }
	defer db.Close()
	
	// We only seed if we want fresh random data. 
	// For now, let's assume we always simulate a new day.
	database.SeedData(db)

	roster, err := scheduler.RunSmartTetris(db)
	if err != nil { log.Fatal(err) }

	// The Goal: Deliver the CSV
	err = scheduler.ExportToCSV(roster, "roster.csv")
	if err != nil { log.Fatal(err) }
}
