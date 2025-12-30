package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/iannsp/shiftopt/internal/database"
	"github.com/iannsp/shiftopt/internal/models"
	"github.com/iannsp/shiftopt/internal/scheduler"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.SeedData(db)

	fmt.Println("========================================")
	fmt.Println("   SHIFTOPT DIAGNOSTIC SUMMARY")
	fmt.Println("========================================")

	// 0. The Context (Demand + Supply)
	printDemandCurve(db)
	printCrewStats(db)

	// 1. Run Strategy A
	rosterHourly, err := scheduler.RunSafeSchedule(db)
	if err != nil {
		log.Fatal(err)
	}
	printStats("Strategy A: Hourly (Fragmented)", rosterHourly)

	// 2. Run Strategy B
	rosterTetris, err := scheduler.RunTetrisSchedule(db)
	if err != nil {
		log.Fatal(err)
	}
	printStats("Strategy B: Tetris (Continuous)", rosterTetris)

	// 3. Comparison
	diff := rosterTetris.TotalCost - rosterHourly.TotalCost
	percent := 0.0
	if rosterHourly.TotalCost > 0 {
		percent = (diff / rosterHourly.TotalCost) * 100
	}

	fmt.Println("\n----------------------------------------")
	fmt.Println("   OPERATIONAL IMPACT ANALYSIS")
	fmt.Println("----------------------------------------")
	if diff > 0 {
		fmt.Printf("Cost of Continuity: +$%.2f (+%.1f%%)\n", diff, percent)
		fmt.Println(">> This is the 'Premium' we pay to give staff 4-hour blocks.")
	} else {
		fmt.Printf("Cost Difference: $%.2f\n", diff)
		fmt.Println(">> Tetris was more efficient. Check for unfilled shifts.")
	}
}

func printDemandCurve(db *sql.DB) {
	fmt.Println("\n[Demand Curve (The Problem)]")
	rows, _ := db.Query("SELECT hour_of_day, needed FROM demands ORDER BY hour_of_day")
	defer rows.Close()

	totalNeeded := 0
	for rows.Next() {
		var h, n int
		rows.Scan(&h, &n)
		totalNeeded += n

		// ASCII Bar Chart
		bar := strings.Repeat("█", n) 
		// If font doesn't support block, use "*"
		// bar := strings.Repeat("*", n) 

		fmt.Printf("  %02d:00 | %s (%d)\n", h, bar, n)
	}
	fmt.Printf("  Total Man-Hours Required: %d\n", totalNeeded)
}

func printCrewStats(db *sql.DB) {
	fmt.Println("\n[Workforce Profile (The Supply)]")
	
	var total, seniors, juniors int
	var avgCost float64

	db.QueryRow("SELECT COUNT(*) FROM employees").Scan(&total)
	db.QueryRow("SELECT COUNT(*) FROM employees WHERE skill_level >= 2").Scan(&seniors)
	db.QueryRow("SELECT COUNT(*) FROM employees WHERE skill_level = 1").Scan(&juniors)
	db.QueryRow("SELECT AVG(hourly_rate) FROM employees").Scan(&avgCost)

	fmt.Printf("  Total Headcount: %d\n", total)
	fmt.Printf("  Composition:     %d Seniors / %d Juniors\n", seniors, juniors)
	fmt.Printf("  Avg Hourly Rate: $%.2f/hr\n", avgCost)
	
	if seniors == 0 {
		fmt.Println("  WARNING: No Seniors in pool! Safety constraints will fail.")
	}
}

func printStats(name string, roster *models.Roster) {
	fmt.Printf("\n[%s]\n", name)
	fmt.Printf("  Total Cost:      $%.2f\n", roster.TotalCost)
	
	assigned := len(roster.Assignments)
	totalNeeded := assigned + roster.Unfilled
	coveragePct := 0.0
	if totalNeeded > 0 {
		coveragePct = (float64(assigned) / float64(totalNeeded)) * 100
	}

	fmt.Printf("  Shifts Assigned: %d / %d (%.1f%%)\n", assigned, totalNeeded, coveragePct)
	
	uniqueStaff := make(map[int]bool)
	for _, a := range roster.Assignments {
		uniqueStaff[a.Employee.ID] = true
	}
	fmt.Printf("  Staff Utilized:  %d people\n", len(uniqueStaff))

	if roster.Unfilled > 0 {
		fmt.Printf("  CRITICAL:        %d Unfilled Shifts\n", roster.Unfilled)
	}

	fmt.Println("  Staff Hours:")
	hoursPerPerson := make(map[string]int)
	var names []string

	for _, a := range roster.Assignments {
		if _, exists := hoursPerPerson[a.Employee.Name]; !exists {
			names = append(names, a.Employee.Name)
		}
		hoursPerPerson[a.Employee.Name]++
	}
	sort.Strings(names) 

	for _, name := range names {
		hours := hoursPerPerson[name]
		fmt.Printf("    - %-15s: %d hrs\n", name, hours)
	}
}
