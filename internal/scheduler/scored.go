package scheduler

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/iannsp/shiftopt/internal/models"
)

// RunSmartTetris uses Penalty Scoring to optimize skill usage AND respects Availability
func RunSmartTetris(db *sql.DB) (*models.Roster, error) {
	fmt.Println("\n--- Generating Smart Tetris Schedule (Penalty Scoring) ---")

	// 1. Fetch Employees
	rows, _ := db.Query("SELECT id, name, hourly_rate, skill_level FROM employees")
	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		rows.Scan(&e.ID, &e.Name, &e.HourlyRate, &e.SkillLevel)
		employees = append(employees, e)
	}
	rows.Close()

	// 1.5 Fetch Unavailability (The Missing Link)
	// Map: EmployeeID -> Map[Hour] -> IsBlocked
	blocked := make(map[int]map[int]bool)
	uRows, err := db.Query("SELECT employee_id, start_hour, end_hour FROM unavailability")
	if err == nil {
		for uRows.Next() {
			var empID, start, end int
			uRows.Scan(&empID, &start, &end)
			
			if blocked[empID] == nil {
				blocked[empID] = make(map[int]bool)
			}
			// Block every hour in the range [start, end)
			for h := start; h < end; h++ {
				blocked[empID][h] = true
			}
		}
		uRows.Close()
	}

	// 2. Fetch Demands
	dRows, _ := db.Query("SELECT hour_of_day, needed FROM demands ORDER BY hour_of_day")
	var sortedHours []int
	demands := make(map[int]int)
	for dRows.Next() {
		var h, n int
		dRows.Scan(&h, &n)
		demands[h] = n
		sortedHours = append(sortedHours, h)
	}
	dRows.Close()

	roster := &models.Roster{}
	shiftEnd := make(map[int]int)
	hoursWorkedTotal := make(map[int]int)

	const MinBlock = 4
	const MaxDaily = 8
	const (
		PenaltySafetyMissing = 1000.0
		PenaltySeniorWaste   = 50.0
	)

	// 3. The Loop
	for _, hour := range sortedHours {
		needed := demands[hour]

		// A. Analyze Current State
		activeCount := 0
		seniorPresent := false
		activeStaff := make(map[int]bool)

		for _, emp := range employees {
			if shiftEnd[emp.ID] > hour {
				// Already working
				isSenior := (emp.SkillLevel >= 2)
				roster.Assignments = append(roster.Assignments, models.Assignment{
					Hour: hour, Employee: emp, IsSenior: isSenior,
				})
				roster.TotalCost += emp.HourlyRate
				hoursWorkedTotal[emp.ID]++
				
				activeCount++
				activeStaff[emp.ID] = true
				if isSenior {
					seniorPresent = true
				}
			}
		}

		// B. Spawn Blocks
		deficit := needed - activeCount
		if deficit > 0 {
			for i := 0; i < deficit; i++ {
				
				type Candidate struct {
					Emp   models.Employee
					Score float64
				}
				var candidates []Candidate

				for _, emp := range employees {
					// --- HARD CONSTRAINTS ---
					
					// 1. Is already working?
					if activeStaff[emp.ID] { continue }
					
					// 2. Will bust 8-hour limit?
					if hoursWorkedTotal[emp.ID]+MinBlock > MaxDaily { continue }

					// 3. **AVAILABILITY CHECK** (The Fix)
					// Check if ANY hour in the proposed block (hour -> hour+4) is blocked
					isBlocked := false
					for b := 0; b < MinBlock; b++ {
						// Logic: If blocked[Alice][09:00] is true, she cannot take a shift starting at 09:00
						// We check hour, hour+1, hour+2, hour+3
						if blocked[emp.ID][hour+b] {
							isBlocked = true
							break
						}
					}
					if isBlocked { continue }


					// --- SOFT CONSTRAINTS (SCORING) ---
					score := emp.HourlyRate

					if !seniorPresent {
						if emp.SkillLevel < 2 {
							score += PenaltySafetyMissing
						}
					} else {
						if emp.SkillLevel >= 2 {
							score += PenaltySeniorWaste
						}
					}

					candidates = append(candidates, Candidate{Emp: emp, Score: score})
				}

				// Sort and Assign
				sort.Slice(candidates, func(i, j int) bool {
					return candidates[i].Score < candidates[j].Score
				})

				if len(candidates) > 0 {
					winner := candidates[0].Emp
					shiftEnd[winner.ID] = hour + MinBlock
					
					isSenior := (winner.SkillLevel >= 2)
					roster.Assignments = append(roster.Assignments, models.Assignment{
						Hour: hour, Employee: winner, IsSenior: isSenior,
					})
					roster.TotalCost += winner.HourlyRate
					hoursWorkedTotal[winner.ID]++
					activeStaff[winner.ID] = true
					if isSenior {
						seniorPresent = true
					}
				} else {
					roster.Unfilled++
				}
			}
		}
	}

	return roster, nil
}
