package main

import (
	"fmt"
	"log"

	"github.com/iannsp/shiftopt/internal/database"
	"github.com/iannsp/shiftopt/internal/scheduler"
)

func main() {
	db, err := database.InitDB()
	if err != nil { log.Fatal(err) }
	defer db.Close()

	// Note: In a real app, we might check existing data, 
	// but here we run the simulation to check the "Scenario"
	database.SeedData(db)

	roster, err := scheduler.RunSafeSchedule(db)
	if err != nil { log.Fatal(err) }

	// The Goal: Show Totals
	fmt.Println("--- Simulation Summary ---")
	fmt.Printf("Total Daily Cost: $%.2f\n", roster.TotalCost)
	fmt.Printf("Total Shifts Assigned: %d\n", len(roster.Assignments))
	
	if roster.Unfilled > 0 {
		fmt.Printf("CRITICAL ALERT: %d shifts were left unfilled due to constraints.\n", roster.Unfilled)
	} else {
		fmt.Println("Status: All shifts covered safely.")
	}
}
