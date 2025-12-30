package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/iannsp/shiftopt/internal/models"
	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "shiftopt.db")
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS employees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		hourly_rate REAL,
		skill_level INTEGER
	);
	CREATE TABLE IF NOT EXISTS demands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hour_of_day INTEGER,
		needed INTEGER
	);
	`
	_, err = db.Exec(schema)
	return db, err
}

func SeedData(db *sql.DB) {
	// 1. Initialize the Random Source based on current time
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	db.Exec("DELETE FROM employees; DELETE FROM demands;")

	// 2. Employees (We keep this pool stable for now, representing "Fixed Staff")
	employees := []models.Employee{
		{Name: "Alice (Vet)", HourlyRate: 50.0, SkillLevel: 2},
		{Name: "Bob (Vet)", HourlyRate: 55.0, SkillLevel: 2},
		{Name: "Carol (Vet)", HourlyRate: 52.0, SkillLevel: 2}, // Added one more Senior
		{Name: "Dave (Jun)", HourlyRate: 20.0, SkillLevel: 1},
		{Name: "Eve (Jun)", HourlyRate: 22.0, SkillLevel: 1},
		{Name: "Frank (Jun)", HourlyRate: 21.0, SkillLevel: 1},
		{Name: "Grace (Grinder)", HourlyRate: 30.0, SkillLevel: 1},
		{Name: "Hank (Grinder)", HourlyRate: 32.0, SkillLevel: 1},
	}

	for _, e := range employees {
		db.Exec("INSERT INTO employees (name, hourly_rate, skill_level) VALUES (?, ?, ?)", e.Name, e.HourlyRate, e.SkillLevel)
	}

	// 3. Generate Randomized Demand (08:00 to 20:00)
	// Logic: Base Curve (Lunch Peak) + Random Noise
	fmt.Println("Seeding Randomized Demand Curve (Sine + Noise)...")
	
	for h := 8; h <= 20; h++ {
		baseNeeded := 2

		// The "Lunch Rush" Pattern
		if h >= 11 && h <= 14 {
			baseNeeded = 5
		}
		// The "Dinner Rush" Pattern
		if h >= 18 && h <= 20 {
			baseNeeded = 4
		}

		// Inject Noise: Randomly add -1 to +2 staff needed
		// This simulates unexpected busloads of customers or quiet days
		noise := rng.Intn(4) - 1 // Generates: -1, 0, 1, or 2
		
		finalNeeded := baseNeeded + noise
		
		// Safety floor: Always need at least 1 person
		if finalNeeded < 1 {
			finalNeeded = 1
		}

		db.Exec("INSERT INTO demands (hour_of_day, needed) VALUES (?, ?)", h, finalNeeded)
	}
}
