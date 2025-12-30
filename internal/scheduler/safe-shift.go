package scheduler

import (
	"database/sql"
	"sort"

	"github.com/iannsp/shiftopt/internal/models"
)

// RunSafeSchedule returns a Roster object instead of printing
func RunSafeSchedule(db *sql.DB) (*models.Roster, error) {
	// 1. Fetch & Sort Employees
	rows, _ := db.Query("SELECT id, name, hourly_rate, skill_level FROM employees")
	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		rows.Scan(&e.ID, &e.Name, &e.HourlyRate, &e.SkillLevel)
		employees = append(employees, e)
	}
	rows.Close()

	sort.Slice(employees, func(i, j int) bool {
		return employees[i].HourlyRate < employees[j].HourlyRate
	})

	// 2. Fetch Demands
	dRows, _ := db.Query("SELECT hour_of_day, needed FROM demands ORDER BY hour_of_day")
	
	roster := &models.Roster{
		Assignments: []models.Assignment{},
	}
	
	hoursWorked := make(map[int]int)
	const MaxDailyHours = 8

	for dRows.Next() {
		var hour, needed int
		dRows.Scan(&hour, &needed)

		assignedThisHour := make(map[int]bool)
		slotsFilled := 0

		// --- PASS 1: Safety (Senior) ---
		for _, emp := range employees {
			if hoursWorked[emp.ID] >= MaxDailyHours { continue }
			if emp.SkillLevel >= 2 {
				hoursWorked[emp.ID]++
				roster.TotalCost += emp.HourlyRate
				roster.Assignments = append(roster.Assignments, models.Assignment{
					Hour: hour, Employee: emp, IsSenior: true,
				})
				assignedThisHour[emp.ID] = true
				slotsFilled++
				break 
			}
		}

		// --- PASS 2: Filler ---
		if slotsFilled < needed {
			for _, emp := range employees {
				if slotsFilled >= needed { break }
				if hoursWorked[emp.ID] >= MaxDailyHours || assignedThisHour[emp.ID] { continue }

				hoursWorked[emp.ID]++
				roster.TotalCost += emp.HourlyRate
				roster.Assignments = append(roster.Assignments, models.Assignment{
					Hour: hour, Employee: emp, IsSenior: false,
				})
				slotsFilled++
			}
		}

		if slotsFilled < needed {
			roster.Unfilled += (needed - slotsFilled)
		}
	}
	dRows.Close()

	return roster, nil
}
