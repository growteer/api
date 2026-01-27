---
name: new-job
description: Add a new Postgres-backed job type + worker handler with retries/backoff and idempotency considerations.
---

# New Job + Worker Handler

When invoked:
1) Define a typed job payload.
2) Add a job creator function in the originating module (command handler).
3) Ensure "state change + enqueue job" can be transactional if needed.
4) Create worker handler:
   - claim job
   - execute
   - handle retries/backoff via nextRunAt + attempts
5) Document idempotency strategy:
   - job-level dedupe key OR
   - safe re-run behavior

Never:
- introduce Rabbit/Kafka
- do fire-and-forget for important jobs
