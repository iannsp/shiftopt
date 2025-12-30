# Engineering Log 009: Optimization Strategy (Penalty Scoring vs. Backtracking)

**Date:** December 30, 2025
**Topic:** Algorithmic Strategy, Heuristics & Performance
**Status:** Architecture Decision Recorded

## Executive Summary

**The Facts:**
*   **The Problem:** Our current "Greedy/Tetris" algorithms are short-sighted. They consume scarce resources (e.g., Seniors) early in the day, causing failures during later high-risk shifts ("Painting ourselves into a corner").
*   **Option A (Backtracking):** A "Brute Force" tree search that tries every combination, undoing decisions when it gets stuck. Guarantees a mathematical solution. Time complexity is Exponential ($O(N!)$).
*   **Option B (Penalty Scoring):** A "Heuristic" approach that adds virtual costs to bad decisions (e.g., "Using a Senior for a junior task costs +$1000"). Time complexity remains Linear/Polynomial ($O(N \log N)$).

**The Conclusion:**
We have decided to implement **Option B: Penalty Scoring**.

---

## Detailed Analysis

### 1. Why we rejected Backtracking (The "Surgery" Case)
Backtracking is the gold standard when "Failure is not an option."
*   **Use Case:** Scheduling Neurosurgeons or ER Nurses.
*   **Why:** If you don't find a valid schedule, people die. You are willing to let the computer run for 4 hours to find the *one* perfect combination that works.
*   **The Retail Reality:** In high-volume retail (our domain), rosters change constantly. Managers need to run "What-If" scenarios in seconds, not hours. A solution that takes 20 minutes to compute because of combinatorial explosion is useless, even if it is "perfect."

### 2. Why we chose Penalty Scoring (The "Manager Intuition")
Penalty Scoring models how a human manager actually thinks. It turns **Soft Constraints** into **Virtual Costs**.

Instead of `Cost = HourlyRate`, the algorithm will see:
$$ \text{Score} = \text{HourlyRate} + (\text{scarcity\_weight} \times \text{penalty}) $$

**The "Senior Preservation" Example:**
*   **Scenario:** 09:00 AM (Quiet Shift).
*   **Alice (Senior):** Rate $50. Penalty for wasting a Senior on a junior task: +$500. **Score: $550.**
*   **Bob (Junior):** Rate $20. Penalty: $0. **Score: $20.**

**Result:** The algorithm naturally picks Bob, "saving" Alice for the evening rush where she is actually needed.

### 3. Operational Benefits
*   **Tunable:** We can adjust behavior without rewriting code. If the business decides "Overtime is worse than Unfilled Shifts," we simply raise the Overtime Penalty weight.
*   **Performance:** It scales efficiently. Scheduling 500 employees takes milliseconds, whereas backtracking would choke.
*   **Explainability:** We can tell a manager: *"The system didn't schedule Alice because we set a high penalty on wasting her skills,"* rather than *"The recursive tree search pruned that branch."*

## Implementation Plan
We will implement a `ScoringEngine` in the next phase.
1.  Define a `Score()` function for an assignment.
2.  Add a `SeniorWastePenalty`.
3.  Add an `OvertimePenalty` (Soft limit before the hard 8-hour limit).


