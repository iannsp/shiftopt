package scheduler

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/iannsp/shiftopt/internal/models"
)


// RunSafeSchedule adds "Max 8h" AND "Min 1 Senior" constraints
func RunSafeSchedule(db *sql.DB) {
	fmt.Println("\n--- Running Safety Net Scheduler (Max 8h + 1 Senior) ---")

	// 1. Fetch & Sort Employees
	rows, _ := db.Query("SELECT id, name, hourly_rate, skill_level FROM employees")
	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		rows.Scan(&e.ID, &e.Name, &e.HourlyRate, &e.SkillLevel)
		employees = append(employees, e)
	}
	rows.Close()

	// Sort by Cost (Cheapest first)
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].HourlyRate < employees[j].HourlyRate
	})

	// 2. Fetch Demands
	dRows, _ := db.Query("SELECT hour_of_day, needed FROM demands ORDER BY hour_of_day")
	
	totalCost := 0.0
	hoursWorked := make(map[int]int)
	const MaxDailyHours = 8
	unassignedShifts := 0

	for dRows.Next() {
		var hour, needed int
		dRows.Scan(&hour, &needed)

		assignedThisHour := make(map[int]bool) // Track who works THIS specific hour
		seniorAssigned := false
		slotsFilled := 0

		// --- PASS 1: The "Safety" Mandate (Find 1 Senior) ---
		for _, emp := range employees {
			// Skip if already working max hours
			if hoursWorked[emp.ID] >= MaxDailyHours {
				continue
			}
			
			// We only want a Senior here
			if emp.SkillLevel >= 2 {
				// Assign the Senior
				hoursWorked[emp.ID]++
				totalCost += emp.HourlyRate
				assignedThisHour[emp.ID] = true
				seniorAssigned = true
				slotsFilled++
				break // We only need ONE senior to satisfy the safety rule
			}
		}

		if !seniorAssigned {
			fmt.Printf("RISK ALERT: Hour %d has NO SENIOR available! (Safety Violation)\n", hour)
		}

		// --- PASS 2: Fill the rest (Cheapest bodies) ---
		if slotsFilled < needed {
			for _, emp := range employees {
				if slotsFilled >= needed {
					break
				}
				// Skip if maxed out OR if already assigned in Pass 1
				if hoursWorked[emp.ID] >= MaxDailyHours || assignedThisHour[emp.ID] {
					continue
				}

				// Assign
				hoursWorked[emp.ID]++
				totalCost += emp.HourlyRate
				slotsFilled++
			}
		}

		if slotsFilled < needed {
			unassignedShifts += (needed - slotsFilled)
		}
	}
	dRows.Close()

	fmt.Printf("Optimization Complete. Daily Cost: $%.2f\n", totalCost)
	if unassignedShifts > 0 {
		fmt.Printf("CRITICAL: %d shifts unfilled.\n", unassignedShifts)
	}
}

