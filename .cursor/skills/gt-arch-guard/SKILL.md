---
name: arch-guard
description: Review a planned change or diff for compliance with the project's architecture decisions (GraphQL-first, ActivityPub edge, Postgres+Prisma, jobs table worker, no Redis/Rabbit in v1).
---

# Architecture Guard

When invoked, do the following:

1) Identify which surface(s) the change touches:
   - public-api (GraphQL)
   - admin-api (GraphQL)
   - federation-edge (ActivityPub HTTP)
   - workers (jobs)
   - shared packages/modules

2) Run a compliance checklist:
   - ✅ Product API remains GraphQL-first (no REST endpoints added for product use-cases)
   - ✅ ActivityPub work stays in federation-edge + workers
   - ✅ Async work uses jobs table + worker (no fire-and-forget for important work)
   - ✅ No Redis/Rabbit/Kafka introduced unless requested
   - ✅ Prisma/Postgres used; module boundaries respected (no cross-module table writes)
   - ✅ Resolvers/controllers are thin; business logic is in commands/queries

3) If violations exist:
   - propose the smallest refactor to become compliant
   - list exact file-level changes to make

Output format:
- “✅ Compliant” OR “❌ Not compliant”
- Bullet list of findings
- Suggested patch plan
