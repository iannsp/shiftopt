# Engineering Log 001: Genesis & Philosophy

**Date:** December 29, 2025
**Team:** Ivo (Product Engineer) & AI Pair
**Status:** Foundation

## 1. Context & The "Why"
This project, **ShiftOpt**, was born from a specific intersection of skills.

The primary developer brings a background in **Operations** and **Product Management**. We are working together to build this to solve the **Mathematical & Operational friction** of workforce scheduling.

We decided to work as a **Pair**:
*   **Ivo:** Provides the business intuition, operational constraints, and architectural direction ("Evolutionary Design").
*   **AI:** Provides syntax generation, pattern matching, and refactoring strategies.

## 2. The Philosophy: "Evolutionary Trees"
We explicitly rejected the idea of "Perfect Code First." 
Instead, we adopted an **Evolutionary Git Strategy**:
*   **Local Branches:** Allowed to be messy, experimental, and chaotic (The "Kitchen").
*   **Main Branch:** Squashed, atomic, and clean (The "Dish").
*   **Code Structure:** We started with a simple script and refactored to `cmd/` + `internal/` only when the complexity demanded it.

## 3. Technology Choices (The Stack)

### Why Go?
Go was chosen for two primary reasons:
1.  **Expertise:** It is the primary developer's language of choice.
2.  **Performance:** It excels at CPU-bound tasks. Scheduling is a combinatorial optimization problem, and Go's performance is required for the calculation loops.

### Why SQLite?
We selected SQLite to **avoid configuration overhead** during the initial phase.
*   It allows us to focus on the algorithm rather than infrastructure.
*   It adheres to our evolutionary principle: we start with the simplest database solution and can migrate to PostgreSQL later if the concurrency requirements demand it.

### Why "Simulation First"?
We rejected the idea of manual data entry.
*   **Decision:** Build a Monte Carlo-style generator (`seedData`) before building the UI.
*   **Reasoning:** If the algorithm can't handle 30 days of high-variance, mathematically generated noise, it won't handle real life.

## 4. Current Architecture
```text
[CLI Entry] -> [Simulation Layer] -> [SQLite DB] -> [Greedy Scheduler] -> [Output]
