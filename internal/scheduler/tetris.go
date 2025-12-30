package scheduler

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/iannsp/shiftopt/internal/models"
)

// RunTetrisSchedule implements Block Scheduling (Min 4 hours contiguous)
func RunTetrisSchedule(db *sql.DB) (*models.Roster, error) {
	fmt.Println("\n--- Generating Tetris Schedule (Block Continuity) ---")

	// 1. Setup Data
	rows, _ := db.Query("SELECT id, name, hourly_rate, skill_level FROM employees")
	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		rows.Scan(&e.ID, &e.Name, &e.HourlyRate, &e.SkillLevel)
		employees = append(employees, e)
	}
	rows.Close()

	// Sort cheapest first
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].HourlyRate < employees[j].HourlyRate
	})

	dRows, _ := db.Query("SELECT hour_of_day, needed FROM demands ORDER BY hour_of_day")
	demands := make(map[int]int)
	var sortedHours []int
	for dRows.Next() {
		var h, n int
		dRows.Scan(&h, &n)
		demands[h] = n
		sortedHours = append(sortedHours, h)
	}
	dRows.Close()

	roster := &models.Roster{}
	
	// Track state
	// shiftEnd[EmployeeID] = The hour their current shift ends (e.g., if set to 14, they work until 14:00)
	shiftEnd := make(map[int]int)
	
	// hoursWorkedTotal[EmployeeID] = Total hours accumulated
	hoursWorkedTotal := make(map[int]int)

	const MinBlock = 4
	const MaxDaily = 8

	// 2. The Tetris Loop
	for _, hour := range sortedHours {
		needed := demands[hour]
		
		// A. Who is ALREADY here? (The Continuity Check)
		activeCount := 0
		activeStaff := make(map[int]bool)

		for _, emp := range employees {
			if shiftEnd[emp.ID] > hour {
				// They are already committed to this block!
				// We MUST assign them (Sunk Cost), even if we don't need them.
				roster.Assignments = append(roster.Assignments, models.Assignment{
					Hour: hour, Employee: emp, IsSenior: false, // Ignoring senior check for MVP
				})
				roster.TotalCost += emp.HourlyRate
				hoursWorkedTotal[emp.ID]++
				activeCount++
				activeStaff[emp.ID] = true
			}
		}

		// B. Do we need MORE people? (Spawn new Blocks)
		deficit := needed - activeCount
		
		if deficit > 0 {
			for i := 0; i < deficit; i++ {
				// Find a fresh person to start a NEW BLOCK
				assigned := false
				
				for _, emp := range employees {
					// 1. Are they already working?
					if activeStaff[emp.ID] { continue }
					
					// 2. Can they take a 4-hour block without busting 8 hours?
					// (Simple check: Just checking total cap for now)
					if hoursWorkedTotal[emp.ID] + MinBlock > MaxDaily { continue }

					// 3. Assign the Block
					// Their shift will end at hour + MinBlock
					shiftEnd[emp.ID] = hour + MinBlock
					
					// Record THIS hour
					roster.Assignments = append(roster.Assignments, models.Assignment{
						Hour: hour, Employee: emp, IsSenior: false,
					})
					roster.TotalCost += emp.HourlyRate
					hoursWorkedTotal[emp.ID]++
					
					activeStaff[emp.ID] = true
					assigned = true
					break
				}
				
				if !assigned {
					roster.Unfilled++
				}
			}
		}
	}

	return roster, nil
}
