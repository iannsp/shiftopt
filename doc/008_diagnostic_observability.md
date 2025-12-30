# Engineering Log 008: Diagnostic Observability

**Date:** December 30, 2025
**Topic:** Tooling, CLI Visualization, and Metrics
**Status:** MVP Validation

## Moving Beyond "Black Boxes"
A complex algorithm is dangerous if you cannot visualize its inputs and outputs. We upgraded our `shiftsummary` tool to provide a complete **Operational Diagnostic**.

## 1. Visualizing the Problem (Demand)
We implemented an ASCII-based Demand Curve directly in the terminal.
```text
09:00 | ██ (2)
12:00 | ██████ (6) <-- The Peak
15:00 | ██ (2)


A full output can be see here:
```bash
Building ShiftOpt (CSV Generator)...
Building ShiftSummary (Diagnostic Tool)...
Build Complete. Artifacts in bin/
Seeding Randomized Demand Curve (Sine + Noise)...
========================================
   SHIFTOPT DIAGNOSTIC SUMMARY
========================================

[Demand Curve (The Problem)]
  08:00 | █ (1)
  09:00 | █ (1)
  10:00 | ███ (3)
  11:00 | ████ (4)
  12:00 | ████ (4)
  13:00 | ███████ (7)
  14:00 | ████ (4)
  15:00 | █ (1)
  16:00 | ████ (4)
  17:00 | ████ (4)
  18:00 | ██████ (6)
  19:00 | ██████ (6)
  20:00 | ██████ (6)
  Total Man-Hours Required: 51

[Workforce Profile (The Supply)]
  Total Headcount: 8
  Composition:     3 Seniors / 5 Juniors
  Avg Hourly Rate: $35.25/hr

[Strategy A: Hourly (Fragmented)]
  Total Cost:      $1574.00
  Shifts Assigned: 48 / 51 (94.1%)
  Staff Utilized:  8 people
  CRITICAL:        3 Unfilled Shifts
  Staff Hours:
    - Alice (Vet)    : 8 hrs
    - Bob (Vet)      : 2 hrs
    - Carol (Vet)    : 6 hrs
    - Dave (Jun)     : 8 hrs
    - Eve (Jun)      : 8 hrs
    - Frank (Jun)    : 8 hrs
    - Grace (Grinder): 4 hrs
    - Hank (Grinder) : 4 hrs

--- Generating Tetris Schedule (Block Continuity) ---

[Strategy B: Tetris (Continuous)]
  Total Cost:      $1769.00
  Shifts Assigned: 55 / 55 (100.0%)
  Staff Utilized:  8 people
  Staff Hours:
    - Alice (Vet)    : 7 hrs
    - Bob (Vet)      : 1 hrs
    - Carol (Vet)    : 7 hrs
    - Dave (Jun)     : 8 hrs
    - Eve (Jun)      : 8 hrs
    - Frank (Jun)    : 8 hrs
    - Grace (Grinder): 8 hrs
    - Hank (Grinder) : 8 hrs

----------------------------------------
   OPERATIONAL IMPACT ANALYSIS
----------------------------------------
Cost of Continuity: +$195.00 (+12.4%)
>> This is the 'Premium' we pay to give staff 4-hour blocks.
```
