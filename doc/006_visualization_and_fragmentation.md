# Engineering Log 006: The Fragmentation Discovery

**Date:** December 30, 2025
**Topic:** Visualization, Continuity & The "Swiss Cheese" Schedule
**Status:** Problem Identified

## 1. The Evidence
We executed our strategy from Log 005 (Visualization First) and generated `roster.csv`.
Upon inspection, the "Fragmentation" hypothesis was confirmed.

**Sample Output (Reconstructed):**
| Hour | Employee | Role | Note |
| :--- | :--- | :--- | :--- |
| 09:00 | Alice | Senior | Working |
| 10:00 | Bob | Senior | Alice replaced by Bob (Cheaper?) |
| 11:00 | Alice | Senior | Alice returns (Demand spiked) |

**The Operational Reality:**
In this scenario, Alice is expected to:
1.  Commute to work for 09:00.
2.  Clock out and wait in the breakroom for an hour at 10:00.
3.  Clock back in at 11:00.

No human being accepts this schedule. This is what we call a **"Swiss Cheese Schedule"**—full of holes.

## 2. The Root Cause
Our current Algorithm (`RunSafeSchedule`) is **Stateless regarding Continuity**.
*   It treats every hour as an independent auction.
*   It asks: *"Who is the cheapest valid person for Hour X?"*
*   It **never** asks: *"Is this person already here?"*

Because we optimized strictly for **Direct Cost (Hourly Rate)**, we ignored **Indirect Cost (Switching Costs)**.

## 3. Solution Brainstorming (Alternative Paths)

To fix this, we need to change the fundamental unit of scheduling from "Hours" to "Blocks." Here are the three engineering strategies we are considering:

### Option A: The "Tetris" Approach (Block Scheduling)
Instead of filling `Hour 9`, `Hour 10`, `Hour 11`... we pre-define **Shift Blocks**.
*   *Definition:* A "Block" is a 4-hour, 6-hour, or 8-hour contiguous chunk.
*   *Logic:* We assign *Blocks* to employees, not hours.
*   *Pros:* Guarantees perfect continuity. 100% human-friendly.
*   *Cons:* Low flexibility. If demand spikes for just 1 hour, we have to pay for a 4-hour block. **Higher Cost.**

### Option B: The "Sticky" Algorithm (Penalty Functions)
We keep the hourly logic but add a "Start-Up Cost."
*   *Logic:* `RealCost = HourlyRate + (IsNewShift ? $50 : $0)`
*   *Mechanism:* If Alice worked Hour 9, she is "cheaper" for Hour 10 than Bob, because Bob has to pay the $50 "commute penalty."
*   *Pros:* Flexible. Mathematically elegant.
*   *Cons:* Hard to tune. If the penalty is too low, fragmentation remains. If too high, costs explode.

### Option C: Post-Processing (The "Defragmenter")
We run the current messy algorithm, then run a second pass to "fill gaps."
*   *Logic:* "If Alice works 9 and 11, force assign her to 10."
*   *Pros:* Easy to implement.
*   *Cons:* Dangerous. Forcing Alice into Hour 10 might violate the "Max 8 Hours" rule we established earlier. It feels like a "hack."

## 4. Recommendation
From a Product Engineering perspective, **Option A (Tetris/Block)** is the most robust solution for Retail/Logistics. It mimics how managers actually think ("I need a morning shifter").

**Option B (Sticky)** is better for Gig Economy apps (Uber/DoorDash), but bad for fixed retail.

We will proceed to investigate implementing **Block-Based Scheduling**.


