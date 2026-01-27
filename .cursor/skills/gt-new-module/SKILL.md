---
name: new-module
description: Scaffold a new bounded-context module (NestJS) following the modular monolith and ports/adapters pattern.
---

# New Module Scaffold (NestJS)

Goal: create a new module that respects boundaries and is ready for future extraction, without adding infrastructure.

Steps:
1) Create module folder with:
   - domain/
   - commands/
   - queries/
   - ports/
   - adapters/ (prisma, http clients, etc.)
2) Export only:
   - Module definition
   - Public ports / use-cases
3) Add unit test skeletons for at least one command and one query.
4) If it requires async work:
   - define job payload type
   - create worker handler stub (do not introduce Rabbit/Kafka)

Deliverables:
- File tree created
- Minimal compile-passing skeleton
- Notes on boundaries and what is exported
