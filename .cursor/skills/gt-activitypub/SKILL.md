---
name: activitypub-implementer
description: Implement ActivityPub edge endpoints and the store-first + async processing pipeline (inbox/outbox + delivery jobs).
---

# ActivityPub Implementer

Rules:
- All HTTP endpoints live in federation-edge.
- Persist raw inbound activities before processing.
- Outgoing federation delivery is always async via jobs.

Workflow:
1) Add/modify controller endpoints in federation-edge.
2) Add persistence records (raw activity + metadata).
3) Enqueue processing/delivery jobs.
4) Implement worker handlers for:
   - processInboxActivity
   - deliverOutboxActivity
5) Ensure idempotency and dedupe.

Output:
- file list
- processing flow diagram (text)
- retry/backoff policy
