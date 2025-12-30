# Engineering Log 003: The Safety Constraint (Cost of Quality)

**Date:** December 30, 2025
**Topic:** Risk Management & Skill Composition
**Status:** Phase 2 (Constraints)

## The Hidden Risk of "Low Cost"
Our previous experiment (the Greedy Algorithm) gave us the lowest financial cost. However, it likely created a dangerous operational reality: **Junior-only shifts.**

If a store is staffed entirely by Level 1 employees (Juniors/Trainees), who handles a payment system crash? Who de-escalates an angry customer? Who closes the safe?

As the Lead Product Engineer put it:
> *"A store run entirely by trainees is cheap but dangerous. We are implementing a Safety Constraint: ensuring operational continuity by mandating leadership presence."*

## The New Rule: "Minimum 1 Senior"
We are introducing a **Composition Constraint**.
*   **Rule:** For every active hour, at least one assigned employee must have `SkillLevel >= 2`.
*   **The Conflict:** Seniors are expensive ($50/hr) compared to Juniors ($20/hr). The Greedy algorithm naturally avoids them.
*   **The Forcing Function:** We will force the scheduler to "pay up" for safety.

## Technical Implementation
We are not just counting bodies anymore; we are inspecting **attributes**.

```go
// Psuedo-code logic
func assignStaff(hour int, needed int) {
    hasSenior := false
    
    // First pass: Try to find a senior
    for _, emp := range employees {
        if emp.SkillLevel == 2 {
            assign(emp)
            hasSenior = true
            break
        }
    }
    
    // Second pass: Fill the rest with cheapest available
    fillRemaining(needed - 1)
}
```

