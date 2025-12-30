package scheduler

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/iannsp/shiftopt/internal/models"
)

// RunConstrained adds the "8-Hour Limit" rule
func RunConstrained(db *sql.DB) {
	fmt.Println("\n--- Running Constrained Scheduler (Max 8h/day) ---")

	// 1. Fetch & Sort Employees (Same as Greedy)
	rows, _ := db.Query("SELECT id, name, hourly_rate FROM employees")
	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		rows.Scan(&e.ID, &e.Name, &e.HourlyRate)
		employees = append(employees, e)
	}
	rows.Close()

	sort.Slice(employees, func(i, j int) bool {
		return employees[i].HourlyRate < employees[j].HourlyRate
	})

	// 2. Fetch Demands
	dRows, _ := db.Query("SELECT hour_of_day, needed FROM demands ORDER BY hour_of_day")
	
	totalCost := 0.0
	unassignedShifts := 0

	// --- THE NEW LOGIC: State Tracking ---
	// We need to remember how many hours each person has worked today.
	// Map: EmployeeID -> Count of Hours
	hoursWorked := make(map[int]int)
	const MaxDailyHours = 8

	for dRows.Next() {
		var hour, needed int
		dRows.Scan(&hour, &needed)

		assignedCount := 0

		// Iterate through employees (Cheapest -> Expensive)
		for _, emp := range employees {
			// Stop if we filled this hour
			if assignedCount >= needed {
				break
			}

			// --- CONSTRAINT CHECK ---
			// If this person has already worked 8 hours, SKIP them.
			// The algorithm is forced to look at the next (more expensive) person.
			if hoursWorked[emp.ID] >= MaxDailyHours {
				continue
			}

			// If valid, assign them
			hoursWorked[emp.ID]++
			totalCost += emp.HourlyRate
			assignedCount++
		}

		// Check if we failed to find enough people
		if assignedCount < needed {
			unassignedShifts += (needed - assignedCount)
			fmt.Printf("WARNING: Hour %d is understaffed! (Ran out of eligible workers)\n", hour)
		}
	}
	dRows.Close()

	fmt.Printf("Optimization Complete. Daily Cost: $%.2f\n", totalCost)
	if unassignedShifts > 0 {
		fmt.Printf("CRITICAL: %d shifts could not be filled due to constraints.\n", unassignedShifts)
	}
}
