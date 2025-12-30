package scheduler

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/iannsp/shiftopt/internal/models"
)

func ExportToCSV(roster *models.Roster, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Hour", "Employee Name", "Role", "Hourly Rate", "Is Safety Senior?"})

	for _, a := range roster.Assignments {
		role := "Junior"
		if a.Employee.SkillLevel == 2 {
			role = "Senior"
		}
		isSafety := "No"
		if a.IsSenior {
			isSafety = "YES"
		}

		writer.Write([]string{
			fmt.Sprintf("%02d:00", a.Hour),
			a.Employee.Name,
			role,
			fmt.Sprintf("%.2f", a.Employee.HourlyRate),
			isSafety,
		})
	}
	fmt.Printf("Success: Roster exported to %s\n", filename)
	return nil
}

