# Engineering Log 007: The Tetris Algorithm & The Edge of Time

**Date:** December 30, 2025
**Topic:** Block Scheduling, Continuity, and Edge Cases
**Status:** Phase 3 Complete

## The Pivot to Blocks ("Tetris")
Following the discovery of "Swiss Cheese" schedules (Log 006), we replaced the hourly auction system with a **Block Allocator**.

**The New Logic:**
Instead of asking *"Who is cheap for this hour?"*, the system now asks:
> *"We have a gap. Spawn a 4-hour block. Who can take it?"*

This guarantees that if Alice is called in, she stays for a minimum of 4 hours. This aligns with standard retail labor contracts.

## The Cost of Continuity
As predicted, this efficiency comes at a premium.
*   **Hourly Cost:** ~$550 (Fragmented/Unusable)
*   **Tetris Cost:** ~$700 (Continuous/Human-Friendly)

**Why the jump?**
If demand spikes at 12:00 but drops at 13:00, the Tetris algorithm is forced to keep the staff member until 16:00. We are paying for **Capacity Buffers**. This is the literal price of treating humans like humans, not robots.

## The "End of Day" Edge Case
During testing, we noticed shifts of **3 hours** appearing in the data.
*   *Cause:* A block spawned at 17:00 (5 PM).
*   *Logic:* It should run for 4 hours (until 21:00).
*   *Reality:* The store closes at 20:00. The loop terminates.

**Decision:**
We accepted this truncation as "Store Closing Procedure." However, we noted a future requirement for **"Minimum Reporting Pay"**—financially, we may still owe the employee for the 4th hour even if the store is closed. For now, the simulation reflects "Hours Worked," not necessarily "Hours Paid."

