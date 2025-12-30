package main

import (
	"log"
    "fmt"
	"github.com/iannsp/shiftopt/internal/database"
	"github.com/iannsp/shiftopt/internal/scheduler"
	"github.com/iannsp/shiftopt/internal/ai"
)

func main() {
	db, err := database.InitDB("shiftopt.db")
	if err != nil { log.Fatal(err) }
	defer db.Close()
	
	// We only seed if we want fresh random data. 
	// For now, let's assume we always simulate a new day.
	database.SeedData(db)


// --- STEP 3: SIMULATE USER INPUT (The "Product" Feature) ---
	// "Alice" is our Senior Vet. Let's block her for the morning.
	incomingText := "Alice has a dentist appointment in the morning"
	
	fmt.Printf("\n[Input] SMS Received: %q\n", incomingText)
	
	// 1. Parse
	constraint := ai.ParseConstraint(incomingText)
	fmt.Printf("[AI] Parsed: Who=%s, When=%d:00-%d:00, Why=%s\n", 
		constraint.EmployeeName, constraint.StartHour, constraint.EndHour, constraint.Reason)

	// 2. Resolve ID
	empID, err := database.GetEmployeeIDByName(db, constraint.EmployeeName)
	if err == nil {
		// 3. Save to DB
		database.AddUnavailability(db, empID, constraint.StartHour, constraint.EndHour, constraint.Reason)
		fmt.Println("[DB] Constraint Saved successfully.")
	} else {
		fmt.Println("[Error] Employee not found.")
	}

	roster, err := scheduler.RunSmartTetris(db)
	if err != nil { log.Fatal(err) }

	// The Goal: Deliver the CSV
	err = scheduler.ExportToCSV(roster, "roster.csv")
	if err != nil { log.Fatal(err) }
}
