package models

// Employee: The resource we need to schedule
type Employee struct {
	ID         int
	Name       string
	HourlyRate float64
	SkillLevel int 
}

// Demand: The requirement 
type Demand struct {
	ID        int
	HourOfDay int 
	Needed    int 
}

// SchedulePlan: The result of our calculation
type SchedulePlan struct {
	Hour      int
	Employee  string
	Cost      float64
}


// Assignment represents one person working one hour
type Assignment struct {
	Hour      int
	Employee  Employee
	IsSenior  bool // Tracks if this person was the "Safety" hire
}

// Roster holds the complete plan for the day
type Roster struct {
	Assignments []Assignment
	TotalCost   float64
	Unfilled    int
}

// Unavailability represents a blocked time slot
type Unavailability struct {
	EmployeeName string // The AI identifies the person by name
	StartHour    int
	EndHour      int
	Reason       string
}
