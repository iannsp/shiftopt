# ShiftOpt: Evolutionary Workforce Scheduling Engine

**ShiftOpt** is a high-performance resource allocation engine designed to solve the "Retail Rostering Problem." 

It bridges the gap between **Operational Constraints** (Labor laws, availability, cost budgets) and **Mathematical Optimization** (Cost minimization, coverage maximization).

> **Project Status:** Phase 1 (MVP - Greedy Algorithm & Simulation)

---

## 📖 The Use Case (The Problem)

In Retail, Logistics, and Hospitality, managers spend 10-20 hours a week manually creating shift rosters. This leads to two critical failures:
1.  **Overstaffing:** Paying for staff during quiet hours (High Cost).
2.  **Understaffing:** Losing revenue during peak hours (Lost Opportunity).

**The Challenge:** 
Assigning 50 employees to 30 days of shifts is not a CRUD problem; it is a combinatorial optimization problem with millions of potential permutations.

**The Solution:** 
ShiftOpt automates this by treating the schedule as a math equation. It ingests "Demand Curves" (predicted foot traffic) and "Employee Constraints" to generate a mathematically optimal roster in milliseconds.

---

## 🧠 Mental Model & Architecture

This project is designed as a **Simulation & Optimization Pipeline**.

### 1. The Simulation Layer (Data Seeding)
Instead of relying on static data, ShiftOpt uses a **Monte Carlo-style simulation** to generate:
*   **Demand Curves:** Uses Sine waves + Noise to simulate realistic retail traffic (Pe


### 2. The Optimization Layer (The Engine)
*   **Current State (Baseline):** A **Greedy Algorithm** that iteratively selects the lowest-cost available resource for every open slot.
*   **Future State:** A **Constraint Satisfaction Solver** (or Genetic Algorithm) that optimizes for "Global Cost" rather than "Local Cost," taking into account fatigue, overtime rules, and skill mixing.

### 3. The Tech Stack
*   **Language:** Go (Golang) - Chosen for raw performance and concurrency in calculation loops.
*   **Database:** SQLite - Embedded, low-latency storage for rapid simulation resets.
*   **Build System:** GNU Make.

---

## 🗺️ Roadmap

We are following an evolutionary development path:

- [x] **Phase 1: Foundation (Current)**
    - [x] Domain Modeling (Employees, Demands).
    - [x] Simulation Engine (Generating realistic test data).
    - [x] Baseline "Greedy" Scheduler (Finding the minimum viable roster).
    - [x] Clean Architecture (`cmd/` vs `internal/`).

- [ ] **Phase 2: The "Real" World (Constraints)**
    - [ ] Implement hard constraints (e.g., "Max 8 hours/day", "Must have 1 Manager on site").
    - [ ] Refactor algorithm to handle backtracking or penalty scoring.

- [ ] **Phase 3: AI Integration**
    - [ ] LLM-based parser: Convert unstructured texts ("I can't work Friday") into structured DB constraints.
    - [ ] Demand Prediction: Use external factors (Weather/Holidays) to adjust demand curves.

- [ ] **Phase 4: Operational Dashboard**
    - [ ] HTML/CSS Visualization of the roster vs. the budget.

---

## 🚀 Quick Start

This project uses `Make` for build automation.

### Prerequisites
*   Go 1.21+
*   Make
*   GCC (for SQLite CGO, though we use a pure-Go driver where possible)

### Commands

```bash
# 1. Run the simulation and scheduler immediately
make run

# 2. Build the optimized binary to ./bin/shiftopt
make build

# 3. Clean up the database and binaries to start fresh
make clean

# 4. Run tests
make test


📂 Project Structure
We follow the standard Go project layout:
.
├── cmd/
│   └── shiftopt/    # Main entry point (The Controller)
├── internal/
│   ├── database/    # SQLite setup and Seeding Logic
│   ├── models/      # Domain entities (Employee, Demand)
│   └── scheduler/   # The Optimization Algorithms
├── bin/             # Compiled artifacts (Ignored by Git)
└── Makefile         # Build automation

```

### Topics 

1.  Combinatorial Optimization Problem
2.  Monte Carlo-style simulation
3.  Operational Constraints 
4.  The Roadmap.

