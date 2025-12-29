package scheduler

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/iannsp/shiftopt/internal/models"
)

func RunGreedy(db *sql.DB) {
	fmt.Println("\n--- Running Greedy Scheduler (Internal Pkg) ---")

	// 1. Fetch Employees
	rows, _ := db.Query("SELECT id, name, hourly_rate FROM employees")
	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		rows.Scan(&e.ID, &e.Name, &e.HourlyRate)
		employees = append(employees, e)
	}
	rows.Close()

	// Sort by Cost
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].HourlyRate < employees[j].HourlyRate
	})

	// 2. Schedule
	dRows, _ := db.Query("SELECT hour_of_day, needed FROM demands ORDER BY hour_of_day")
	var plan []models.SchedulePlan
	totalCost := 0.0

	for dRows.Next() {
		var hour, needed int
		dRows.Scan(&hour, &needed)

		if needed > len(employees) {
			continue
		}

		for i := 0; i < needed; i++ {
			emp := employees[i]
			plan = append(plan, models.SchedulePlan{
				Hour:     hour,
				Employee: emp.Name,
				Cost:     emp.HourlyRate,
			})
			totalCost += emp.HourlyRate
		}
	}
	dRows.Close()

	fmt.Printf("Optimization Complete. Daily Cost: $%.2f\n", totalCost)
}
