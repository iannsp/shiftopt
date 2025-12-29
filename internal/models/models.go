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
