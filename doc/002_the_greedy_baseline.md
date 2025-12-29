# Engineering Log 002: The Greedy Baseline

**Date:** December 29, 2025
**Topic:** Algorithms & Baselines
**Status:** Phase 1 Complete

## The Trap of "Premature Optimization"
In workforce scheduling, the mathematical possibilities are endless. Assigning 50 employees to hourly slots over 30 days results in a search space larger than the number of atoms in the universe (Combinatorial Explosion).

The temptation is to immediately reach for complex tools: Linear Programming, Genetic Algorithms, or AI.

We chose to do the opposite. We implemented a **Greedy Algorithm**.

## What is the Greedy Scheduler?
A "Greedy" approach makes the locally optimal choice at each stage.
*   **The Logic:** For Hour 9, who is the cheapest available person? Assign them. Move to Hour 10. Repeat.
*   **The Code (Go):**
    ```go
    // Sort employees by cost (Cheapest -> Most Expensive)
    sort.Slice(employees, func(i, j int) bool {
        return employees[i].HourlyRate < employees[j].HourlyRate
    })

    // Fill the slot
    for i := 0; i < needed; i++ {
        plan = append(plan, employees[i])
    }
    ```

## Why build "Bad" Logic?
We know this algorithm is flawed operationally.
1.  **No Fatigue:** It will assign the cheapest employee (let's call him "Charlie") to work 24 hours a day, 7 days a week.
2.  **No Fairness:** Expensive veterans will never get a shift.

However, in Engineering terms, this is our **Control Variable**.

By running the Greedy Scheduler, we established the **Theoretical Minimum Cost** (The "Floor").
*   *Example Result:* $450/day.

Now, when we build the complex algorithm (that respects labor laws), the cost will inevitably rise (e.g., to $600/day).
The difference ($150) is the **"Cost of Compliance."**

## The Ops Insight
As a former Operations leader, I know that you cannot optimize what you cannot measure.
The Greedy Baseline gives us the ability to tell the business:
> *"Here is the cheapest possible schedule (illegal). Here is the legal schedule. The gap is the price of our constraints."*

## Next Steps
Now that we have the Baseline, Phase 2 is about applying **Constraints**. We will implement rules to stop "Charlie" from working 24 hours straight, forcing the algorithm to pick the second-cheapest option.
