# Engineering Log 002: The Greedy Baseline & Metrics Strategy

**Date:** December 29, 2025
**Topic:** Algorithms, Heuristics & Key Performance Indicators (KPIs)
**Status:** Phase 1 Complete

## The Trap of "Premature Optimization"
In workforce scheduling, the mathematical possibilities are endless. Assigning 50 employees to hourly slots over 30 days results in a search space larger than the number of atoms in the universe (Combinatorial Explosion).

The temptation is to immediately reach for complex tools: Linear Programming, Genetic Algorithms, or AI. We chose to do the opposite. We implemented a **Greedy Algorithm**.

## What is a "Greedy" Approach?
In Computer Science, a "Greedy Algorithm" is a paradigm that builds a solution piece by piece, always choosing the next piece that offers the most immediate benefit. It focuses entirely on **Local Optima** (what is best *right now*) rather than **Global Optima** (what is best for the *whole month*).

Imagine a hiker trying to climb the highest peak of a mountain range in thick fog.
*   **A Greedy Hiker** looks at their feet and takes the steepest step upward available.
*   **The Problem:** They will climb very fast, but they might end up stuck on a small hill, missing the true summit because to get there, they needed to go *down* into a valley first.

In our Scheduler:
1.  The algorithm looks at **Hour 9**.
2.  It sorts employees by price.
3.  It picks "Charlie" ($20/hr) because he is the cheapest.
4.  It moves to **Hour 10** and picks Charlie again.
It never asks: *"Should I save Charlie for the busy evening shift?"* It just consumes the cheapest resource immediately.

## Why build "Bad" Logic? The Power of Baselines.
We know this algorithm is flawed operationally. It will assign "Charlie" to work 24 hours a day, causing burnout and violating labor laws.

However, in Engineering terms, this is our **Control Variable**. We are using the Greedy approach to establish the **Theoretical Minimum** for several key metrics.

### Defines the "Cost of Compliance"
*   **Greedy Cost:** $450/day (Illegal, but cheap).
*   **Legal Cost:** $600/day (Future algorithm).
*   **Insight:** The $150 difference is the "Price of Labor Laws." This allows the business to quantify exactly how much constraints cost them.

### Defines the "Fairness Gap"
We can measure the distribution of shifts among staff (Standard Deviation or Gini Coefficient).
*   **Greedy:** 1 person does 100% of the work (Maximum Unfairness).
*   **Target:** Work is distributed according to availability.
*   **Insight:** We can track how much "extra cost" we incur effectively to "buy" employee satisfaction/fairness.

### Defines "Skill Waste"
*   **Greedy:** Uses a Senior Engineer for a Junior task because they happened to be the cheapest available at that specific second.
*   **Target:** Save Seniors for high-leverage moments.
*   **Insight:** We can measure **"Skill Utilization Rate"**—what % of a Senior's time was spent on Senior tasks?

## Next Steps
Now that we have the Baseline and the Metrics defined, **Phase 2** is about applying **Constraints**.

We will implement rules to stop "Charlie" from working 24 hours straight, effectively forcing the algorithm to "walk down into the valley" (pick a more expensive person now) to reach the "summit" (a valid, legal schedule).

