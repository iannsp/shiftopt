package tests

import (
	"testing"

	"github.com/iannsp/shiftopt/internal/ai"
	"github.com/iannsp/shiftopt/internal/database"
	"github.com/iannsp/shiftopt/internal/scheduler"
)

func TestAIConstraintIntegration(t *testing.T) {
	// 1. Setup In-Memory DB (Fast & Isolated)
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	defer db.Close()

	// 2. Define Scenarios (The Table)
	scenarios := []struct {
		name          string
		inputText     string
		expectedEmp   string
		expectedStart int
		expectedEnd   int // Logic treats this as Exclusive (Start <= h < End)
	}{
		{
			name:          "Alice Morning Dentist",
			inputText:     "Alice has a dentist appointment in the morning",
			expectedEmp:   "Alice (Vet)",
			expectedStart: 8,
			expectedEnd:   12, // Blocks 08, 09, 10, 11
		},
		{
			name:          "Bob Afternoon Sick",
			inputText:     "Bob is sick in the afternoon",
			expectedEmp:   "Bob (Vet)",
			expectedStart: 13,
			expectedEnd:   17, // Blocks 13, 14, 15, 16
		},
		{
			name:          "Carol All Day Personal",
			inputText:     "Carol is away all day",
			expectedEmp:   "Carol (Vet)",
			expectedStart: 8,
			expectedEnd:   20, // Blocks everything
		},
		{
			name:          "Default Fallback",
			inputText:     "Alice is busy",
			expectedEmp:   "Alice (Vet)",
			expectedStart: 9,
			expectedEnd:   10, // Default 1-hour block
		},
	}

	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			// A. Clean State (Reset DB for every scenario)
			database.SeedData(db)

			// B. The AI Step using Mock
			constraint := ai.MockParse(tc.inputText)
            // the AI Step using Gemini
			//constraint := ai.ParseConstraint(tc.inputText)
			
			// Verify AI Parsing
			if constraint.EmployeeName != tc.expectedEmp {
				t.Errorf("AI Parsing Failed. Want Name='%s', Got='%s'", tc.expectedEmp, constraint.EmployeeName)
			}
			if constraint.StartHour != tc.expectedStart || constraint.EndHour != tc.expectedEnd {
				t.Errorf("AI Parsing Failed. Want Time=%d-%d, Got=%d-%d", 
					tc.expectedStart, tc.expectedEnd, constraint.StartHour, constraint.EndHour)
			}

			// C. The DB Step (Save Constraint)
			id, err := database.GetEmployeeIDByName(db, constraint.EmployeeName)
			if err != nil {
				t.Fatalf("Could not find employee '%s' in DB. Check SeedData.", constraint.EmployeeName)
			}
			err = database.AddUnavailability(db, id, constraint.StartHour, constraint.EndHour, constraint.Reason)
			if err != nil {
				t.Fatalf("Failed to save constraint to DB: %v", err)
			}

			// D. The Scheduler Step (Run Smart Tetris)
			roster, err := scheduler.RunSmartTetris(db)
			if err != nil {
				t.Fatalf("Scheduler crashed: %v", err)
			}

			// E. The Verification Step (Check Roster)
			for _, assignment := range roster.Assignments {
				// Only check the employee relevant to this scenario
				if assignment.Employee.Name == tc.expectedEmp {
					// Check if they are working inside the forbidden zone
					if assignment.Hour >= tc.expectedStart && assignment.Hour < tc.expectedEnd {
						t.Errorf("CONSTRAINT VIOLATION: %s was assigned to Hour %d (Blocked Range: %d-%d)", 
							tc.expectedEmp, assignment.Hour, tc.expectedStart, tc.expectedEnd)
					}
				}
			}
		})
	}
}
