package scheduler

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/iannsp/shiftopt/internal/models"
)

// RunSmartTetris uses Penalty Scoring to optimize skill usage
func RunSmartTetris(db *sql.DB) (*models.Roster, error) {
	fmt.Println("\n--- Generating Smart Tetris Schedule (Penalty Scoring) ---")

	// 1. Fetch Data
 blocked := make(map[int]map[int]bool)
    
    uRows, _ := db.Query("SELECT employee_id, start_hour, end_hour FROM unavailability")
    for uRows.Next() {
        var empID, start, end int
        uRows.Scan(&empID, &start, &end)
        
        if blocked[empID] == nil {
            blocked[empID] = make(map[int]bool)
        }
        for h := start; h < end; h++ { // [start, end)
            blocked[empID][h] = true
        }
    }
    uRows.Close()

	rows, _ := db.Query("SELECT id, name, hourly_rate, skill_level FROM employees")
	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		rows.Scan(&e.ID, &e.Name, &e.HourlyRate, &e.SkillLevel)
		employees = append(employees, e)
	}
	rows.Close()

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

	// --- PENALTY CONFIGURATION ---
	const (
		PenaltySafetyMissing = 1000.0 // Huge penalty if we pick a Junior when we need a Senior
		PenaltySeniorWaste   = 50.0   // Medium penalty if we pick a Senior when we don't need one
	)

	// 2. The Loop
	for _, hour := range sortedHours {
		needed := demands[hour]

		// A. Analyze Current State
		activeCount := 0
		seniorPresent := false
		activeStaff := make(map[int]bool)

		for _, emp := range employees {
			if shiftEnd[emp.ID] > hour {
				// This person is working this hour
				isSenior := (emp.SkillLevel >= 2)
             // 1. Hard Constraints
            if activeStaff[emp.ID] { continue }
            if hoursWorkedTotal[emp.ID]+MinBlock > MaxDaily { continue }
            isBlocked := false
            for b := 0; b < MinBlock; b++ {
                if blocked[emp.ID][hour+b] {
                    isBlocked = true
                    break
                }
            }
            if isBlocked { continue } 

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

		// B. Spawn Blocks based on Score
		deficit := needed - activeCount
		if deficit > 0 {
			for i := 0; i < deficit; i++ {
				
				// --- THE SCORING ENGINE ---
				// We re-evaluate candidates for THIS specific slot
				type Candidate struct {
					Emp   models.Employee
					Score float64
				}
				var candidates []Candidate

				for _, emp := range employees {
					// 1. Hard Constraints (Availability/Max Hours)
					if activeStaff[emp.ID] { continue }
					if hoursWorkedTotal[emp.ID]+MinBlock > MaxDaily { continue }

					// 2. Calculate Score
					score := emp.HourlyRate // Start with Base Cost

					if !seniorPresent {
						// Context: We URGENTLY need a Senior for safety
						if emp.SkillLevel < 2 {
							score += PenaltySafetyMissing // "Don't pick this Junior!"
						}
						// Seniors get no penalty (score = base rate), so they win
					} else {
						// Context: Safety is already handled. We just need bodies.
						if emp.SkillLevel >= 2 {
							score += PenaltySeniorWaste // "Save this Senior for later!"
						}
						// Juniors get no penalty, so they win
					}

					candidates = append(candidates, Candidate{Emp: emp, Score: score})
				}

				// 3. Sort by Score (Lowest wins)
				sort.Slice(candidates, func(i, j int) bool {
					return candidates[i].Score < candidates[j].Score
				})

				// 4. Assign the Winner
				if len(candidates) > 0 {
					winner := candidates[0].Emp
					
					shiftEnd[winner.ID] = hour + MinBlock
					
					// Record assignments
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
