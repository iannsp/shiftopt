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
	if err != nil { log.Fatal(err) }
	defer db.Close()
	database.SeedData(db)

	fmt.Println("========================================")
	fmt.Println("   SHIFTOPT DIAGNOSTIC SUMMARY")
	fmt.Println("========================================")

	// 1. Context: The Workforce
	printCrewStats(db)

	// 2. Execution: Run All 3 Strategies
	// A. Baseline
	rosterHourly, err := scheduler.RunSafeSchedule(db)
	if err != nil { log.Fatal(err) }

	// B. Block Logic (Dumb)
	rosterTetris, err := scheduler.RunTetrisSchedule(db)
	if err != nil { log.Fatal(err) }

	// C. Scored Logic (Smart)
	rosterSmart, err := scheduler.RunSmartTetris(db)
	if err != nil { log.Fatal(err) }

	// 3. Visualization: Inspect the "Smart" Roster deeply
	printVisualDistribution(db, rosterSmart)

	// 4. Comparison: The Numbers
	fmt.Println("\n[Strategy Showdown: Cost vs. Coverage]")
	printSummaryRow("1. Hourly (Fragmented)", rosterHourly)
	printSummaryRow("2. Tetris (Basic Block)", rosterTetris)
	printSummaryRow("3. Smart  (Scored Block)", rosterSmart)

	// 5. The "Smart" Delta Analysis
	diff := rosterTetris.TotalCost - rosterSmart.TotalCost
	fmt.Println("\n[Optimization Analysis]")
	if diff > 0 {
		fmt.Printf(">> SUCCESS: Scoring Engine saved $%.2f compared to basic Tetris.\n", diff)
		fmt.Println("   (Optimized usage of Seniors during low-risk hours)")
	} else if diff < 0 {
		fmt.Printf(">> NOTE: Smart Engine cost $%.2f MORE than basic Tetris.\n", -diff)
		fmt.Println("   (Likely forced expensive Seniors to cover Safety gaps that Basic missed)")
	} else {
		fmt.Println(">> NEUTRAL: No cost difference. Workforce Constraints were loose.")
	}
}

// --- VISUALIZATION HELPERS ---

func printVisualDistribution(db *sql.DB, roster *models.Roster) {
	fmt.Println("\n[Smart Schedule Composition]")
	fmt.Println("Legend: [V]eteran, [J]unior, [G]rinder, [_]Missed")

	// Get Demands
	demands := make(map[int]int)
	rows, _ := db.Query("SELECT hour_of_day, needed FROM demands")
	var hours []int
	for rows.Next() {
		var h, n int
		rows.Scan(&h, &n)
		demands[h] = n
		hours = append(hours, h)
	}
	rows.Close()
	sort.Ints(hours)

	// Map Roster
	allocations := make(map[int][]string)
	for _, a := range roster.Assignments {
		char := "J"
		if strings.Contains(a.Employee.Name, "(Vet)") {
			char = "V"
		} else if strings.Contains(a.Employee.Name, "(Grinder)") {
			char = "G"
		}
		allocations[a.Hour] = append(allocations[a.Hour], char)
	}

	// Render
	for _, h := range hours {
		needed := demands[h]
		staff := allocations[h]
		
		// Sort: V -> G -> J
		sort.Slice(staff, func(i, j int) bool {
			order := map[string]int{"V": 0, "G": 1, "J": 2}
			return order[staff[i]] < order[staff[j]]
		})

		var barBuilder strings.Builder
		for _, c := range staff {
			barBuilder.WriteString("[" + c + "]")
		}
		
		missing := needed - len(staff)
		if missing > 0 {
			for i := 0; i < missing; i++ {
				barBuilder.WriteString("[_]")
			}
		}

		fmt.Printf("  %02d:00 | %-25s (Target: %d)\n", h, barBuilder.String(), needed)
	}
}

// --- STATS HELPERS ---

func printSummaryRow(label string, r *models.Roster) {
	assigned := len(r.Assignments)
	totalNeeded := assigned + r.Unfilled
	
	// Create a status string
	status := "OK"
	if r.Unfilled > 0 {
		status = fmt.Sprintf("MISSING %d", r.Unfilled)
	}

	fmt.Printf("  %-25s | Cost: $%7.2f | Cov: %d/%d | %s\n", 
		label, r.TotalCost, assigned, totalNeeded, status)
}

func printCrewStats(db *sql.DB) {
	fmt.Println("\n[Workforce Supply]")
	var total, seniors, juniors int
	var avgCost float64
	db.QueryRow("SELECT COUNT(*) FROM employees").Scan(&total)
	db.QueryRow("SELECT COUNT(*) FROM employees WHERE skill_level >= 2").Scan(&seniors)
	db.QueryRow("SELECT COUNT(*) FROM employees WHERE skill_level = 1").Scan(&juniors)
	db.QueryRow("SELECT AVG(hourly_rate) FROM employees").Scan(&avgCost)
	fmt.Printf("  Headcount: %d (%d Seniors, %d Juniors)\n", total, seniors, juniors)
}
