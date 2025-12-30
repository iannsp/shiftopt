package ai

import (
	"strings"
	"github.com/iannsp/shiftopt/internal/models"
)

// MockParse simulates an LLM extracting data from text.
// Input: "Alice has a dentist appointment in the morning"
// Output: {Name: "Alice", Start: 9, End: 12, Reason: "Dentist"}
func MockParse(input string) models.Unavailability {
	input = strings.ToLower(input)
	result := models.Unavailability{}

	// 1. Identify Who (Simple substring match)
	// In a real app, this would be Fuzzy Matching against the DB
	if strings.Contains(input, "alice") {
		result.EmployeeName = "Alice (Vet)"
	} else if strings.Contains(input, "bob") {
		result.EmployeeName = "Bob (Vet)"
	} else if strings.Contains(input, "carol") {
		result.EmployeeName = "Carol (Vet)"
	}

	// 2. Identify When (Heuristics)
	if strings.Contains(input, "morning") {
		result.StartHour = 8
		result.EndHour = 12
	} else if strings.Contains(input, "afternoon") {
		result.StartHour = 13
		result.EndHour = 17
	} else if strings.Contains(input, "all day") {
		result.StartHour = 8
		result.EndHour = 20
	} else {
		// Default to a specific slot if unspecified
		result.StartHour = 9
		result.EndHour = 10
	}

	// 3. Identify Why
	if strings.Contains(input, "dentist") {
		result.Reason = "Medical (Dentist)"
	} else if strings.Contains(input, "sick") {
		result.Reason = "Medical (Sick)"
	} else {
		result.Reason = "Personal"
	}

	return result
}
