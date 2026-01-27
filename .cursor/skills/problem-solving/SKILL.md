---
name: problem-solving
description: Default problem-solving policy. Use for architecture and implementation decisions, including features, bug fixes and refactorings.
---

# Simplicity First

## Operating Principle
Default to the **simplest option that ships safely**:
- least code
- fewest moving parts
- fewest dependencies
- clearest mental model
- easiest to delete later

If two options are both “good enough”, pick the one that is **simpler to reason about** and **cheaper to change**.

---

## Mandatory Pre-Implementation Pause (Do this before writing code)
Ask yourself:

1) **Is there a simpler way?**  
   Can this be done with less code, fewer components, fewer abstractions, fewer deps?

2) **Do I really need this?**  
   If the requirement is unclear, confirm it instead of building “just in case”.

If any requirement is ambiguous, ask **one** targeted question before implementing.

---

## Challenge Your First Instinct (Required)
Before committing to a solution:
- Write down **Option A / B / C** (2–3 viable approaches).
- Pick the simplest option that meets the actual requirement.
- If you choose a more complex option, explicitly justify why the simpler options fail.

---

## Prefer Built-ins Over Custom Code
Use framework/platform features first:
- Use built-in APIs, conventions, and primitives before inventing utilities, hooks, or state layers.
- Only introduce a custom abstraction if it:
  - removes repeated complexity, **and**
  - has a clear “delete path”, **and**
  - you can show at least **two** call sites (or a near-term second one).

---

## Avoid Premature Optimization
Do **not** add complexity for hypothetical futures:
- no “maybe we’ll need it later” layers
- no speculative caching, message buses, CQRS, or global state without a real pressure signal
- no new dependency unless it buys clear, immediate value

If performance is a concern, measure first (or add a minimal probe), then optimize the bottleneck.

---

## Anti-Patterns (Hard “No” unless proven necessary)
- Adding abstractions with a single call site
- Introducing new dependencies “because it’s common”
- Building a generalized framework when a small function solves the task
- Turning simple flows into event-driven architectures without pressure
- Adding caching, batching, queues, or CQRS when a straightforward request/response works
