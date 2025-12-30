# Engineering Log 005: Product Checkpoint & Gap Analysis

**Date:** December 30, 2025
**Topic:** Product Strategy, Gap Analysis & Pivot
**Status:** Reviewing Phase 2

## The Hard Question: "Is this useful?"
We paused development to ask a critical Product Management question: *If we shipped this binary today, would it solve a real user problem?*

**The Verdict:**
*   **As an Operational Tool? NO.**
    *   We cannot give this to a Store Manager. It outputs a "Cost" ($600), not a "Roster" (Alice: 09:00-17:00).
    *   It treats humans like robots: it assigns scattered single hours (09:00, then 11:00, then 14:00) rather than contiguous shifts.
*   **As a Strategic Audit Tool? YES.**
    *   We *can* give this to a COO or Finance Director.
    *   It successfully answers: *"Do we have the physical capacity to staff our stores safely?"*
    *   It quantifies the "Cost of Compliance" and "Cost of Safety."

## The Gap Analysis (What is missing?)
To move from a "Back-Office Simulator" to a "Front-Line Product," we identified three critical gaps:

1.  **The "Fragmentation" Problem (Continuity)**
    *   *Issue:* Real shifts are blocks (e.g., 4-8 hours). Our algorithm currently assigns disconnected hours based on pure cost efficiency.
    *   *Impact:* No human would accept the generated schedule.

2.  **The "Data Entry" Friction (Availability)**
    *   *Issue:* Real employees have complex lives (Dentist appointments, Classes, Childcare).
    *   *Current State:* Entering these constraints manually into a database is "Data Entry Hell." Managers hate it, so they won't use the tool.
    *   *Impact:* The system assumes 100% availability, rendering the output invalid for specific days.

3.  **The "Black Box" Problem (Visualization)**
    *   *Issue:* CLI text output is insufficient for roster verification.
    *   *Impact:* Trust is low without a visual grid.

## Decision Record: Visualization vs. Data Entry
We debated which path to prioritize for Phase 3: **Input (AI)** vs. **Output (Visualization)**.

### Option A: Data Entry First (AI Availability)
*   **Pros:** The user can finally input reality (e.g., "I can't work Friday"). The calculation becomes operationally valid.
*   **Cons:** The user **cannot verify** the result. The system might say "Schedule Generated," but without a visual grid, the manager cannot confirm if the AI correctly blocked the specific slot. It remains a "Black Box."

### Option B: Visualization First (HTML Roster)
*   **Pros:** The user can **see** exactly who is working when. It turns the "Black Box" into a "Glass Box," effectively exposing algorithmic stupidity (like fragmented shifts).
*   **Cons:** The schedule is still "idealized" because it ignores availability constraints.

### The Verdict: Visualization First
We chose **Option B**.
**Reasoning:** *You cannot fix what you cannot see.* We need to turn the lights on (Visualization) to debug the scheduling logic before we add the complexity of AI constraints.

## The Pivot
**Phase 3 Strategy:**
1.  Refactor the Scheduler to return structured data (`Roster` objects) instead of just printing text.
2.  Build a **HTML Generator** to render a 08:00-20:00 grid showing exactly who is working when.

