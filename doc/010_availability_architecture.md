# Engineering Log 010: Availability Architecture (The Anti-Roster)

**Date:** December 30, 2025
**Topic:** Negative Constraints & Unstructured Data
**Status:** Phase 3 Start

## The Problem
Our current scheduler assumes **Infinite Capacity**. It believes every employee is available 24/7.
In reality, availability is the primary constraint in retail.

## The "Anti-Roster"
We are introducing a concept called the "Anti-Roster."
*   **Roster:** Who IS working.
*   **Anti-Roster:** Who CANNOT work.

## Schema Design
We need to store time-blocks.
```sql
CREATE TABLE unavailability (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER,
    start_hour INTEGER, -- e.g., 9
    end_hour INTEGER,   -- e.g., 11 (meaning 09:00 to 11:00)
    reason TEXT,        -- "Dentist", "Class", "Sick"
    FOREIGN KEY(employee_id) REFERENCES employees(id)
);
```
