package database

import (
	"database/sql"
	"fmt"
_	"log"

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
	db.Exec("DELETE FROM employees; DELETE FROM demands;")

	employees := []models.Employee{
		{Name: "Alice (Vet)", HourlyRate: 50.0, SkillLevel: 2},
		{Name: "Bob (Vet)", HourlyRate: 55.0, SkillLevel: 2},
		{Name: "Charlie (Jun)", HourlyRate: 20.0, SkillLevel: 1},
		{Name: "Dave (Jun)", HourlyRate: 22.0, SkillLevel: 1},
		{Name: "Eve (Grinder)", HourlyRate: 30.0, SkillLevel: 1},
	}

	for _, e := range employees {
		db.Exec("INSERT INTO employees (name, hourly_rate, skill_level) VALUES (?, ?, ?)", e.Name, e.HourlyRate, e.SkillLevel)
	}

	hours := []int{9, 10, 11, 12, 13, 14, 15, 16, 17}
	for _, h := range hours {
		needed := 2
		if h == 12 || h == 13 {
			needed = 4 
		}
		db.Exec("INSERT INTO demands (hour_of_day, needed) VALUES (?, ?)", h, needed)
	}
	fmt.Println("Database seeded.")
}
